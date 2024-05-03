package serve

import (
	"github.com/apenella/ransidble/internal/configuration"
	ansibleplaybookservice "github.com/apenella/ransidble/internal/domain/core/service/ansibleplaybook"
	server "github.com/apenella/ransidble/internal/handler/http"
	ansibleplaybookhandler "github.com/apenella/ransidble/internal/handler/http/ansible-playbook"
	ansibleplaybookrunner "github.com/apenella/ransidble/internal/infrastructure/execute/ansible-playbook"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/cobra"
)

func NewCommand(config *configuration.Configuration) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Serve is a command to start a Ransidble server",
		Long:  "Serve is a command to start a Ransidble server",
		RunE: func(cmd *cobra.Command, args []string) error {

			router := echo.New()
			router.Use(middleware.Logger())
			router.Use(middleware.GzipWithConfig(middleware.GzipConfig{
				Level: 5,
			}))

			ansibleplaybookrunner := ansibleplaybookrunner.NewAnsiblePlaybookRun()
			ansiblePlaybookService := ansibleplaybookservice.NewAnsiblePlaybookService(ansibleplaybookrunner)
			ansiblePlaybookHandler := ansibleplaybookhandler.NewAnsiblePlaybookHandler(ansiblePlaybookService)

			router.POST("/command/ansible-playbook", ansiblePlaybookHandler.Handle)

			s := server.NewServer(config.HTTPListenAddress, router)

			err := s.Start()
			if err != nil {
				cmd.Println("Error starting server", err)
				return err
			}
			return nil
		},
	}

	return cmd
}
