package http

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/apenella/ransidble/internal/domain/ports/repository"
	"github.com/labstack/echo/v4"
)

var (
	// ErrServerStarting represents an error when starting the server
	ErrServerStarting = fmt.Errorf("error starting server")

	// ErrServerStopping represents an error when stopping the server
	ErrServerStopping = fmt.Errorf("error stopping server")
)

// Server represents a HTTP server
type Server struct {
	logger   repository.Logger
	once     sync.Once
	server   *http.Server
	stopCh   chan struct{}
	stopOnce sync.Once
}

// NewServer creates a new server
func NewServer(listenAddress string, handler *echo.Echo, logger repository.Logger) *Server {
	return &Server{
		server: &http.Server{
			Addr:    listenAddress,
			Handler: handler,
		},
		stopCh: make(chan struct{}),
		logger: logger,
	}
}

// Start starts the server
func (s *Server) Start(ctx context.Context) (err error) {

	s.once.Do(func() {

		s.logger.Info(fmt.Sprintf("Starting server on %s", s.server.Addr), map[string]interface{}{
			"component": "Server.Start",
			"package":   "github.com/apenella/ransidble/internal/handler/http",
			"address":   s.server.Addr,
		})

		if ctx == nil {
			ctx = context.Background()
		}

		errListenAndServeCh := make(chan error)
		go func() {
			errListenAndServe := s.server.ListenAndServe()
			if errListenAndServe != nil {
				errListenAndServeCh <- errListenAndServe
			}
		}()

		select {
		case errListenAndServe := <-errListenAndServeCh:

			if errListenAndServe != nil {
				err = fmt.Errorf("%w: %s", ErrServerStarting, errListenAndServe)
				s.Stop()
			}
		case <-s.stopCh:

			errShutdown := s.server.Shutdown(ctx)
			if errShutdown != nil {
				err = fmt.Errorf("%w: %s", ErrServerStopping, errShutdown)
			}

			s.logger.Info("HTTP Server stopped", map[string]interface{}{
				"component": "Server.Start",
				"package":   "github.com/apenella/ransidble/internal/handler/http",
			})
			return
		}
	})

	return
}

// Stop stops the server
func (s *Server) Stop() {
	s.logger.Info("Stopping HTTP server", map[string]interface{}{
		"component": "Server.Stop",
		"package":   "github.com/apenella/ransidble/internal/handler/http",
	})

	s.stopOnce.Do(func() {
		close(s.stopCh)
	})
}
