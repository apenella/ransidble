package http

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/labstack/echo/v4"
)

type Server struct {
	server *http.Server
}

func NewServer(listenAddress string, handler *echo.Echo) *Server {
	return &Server{
		server: &http.Server{
			Addr:    listenAddress,
			Handler: handler,
		},
	}
}

func (s *Server) Start() error {
	var err error

	log.Printf("Starting server on %s", s.server.Addr)

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	errListenAndServe := make(chan error)

	go func() {
		errListenAndServe <- s.server.ListenAndServe()
	}()

	select {
	case err = <-errListenAndServe:
		if err != nil {
			return fmt.Errorf("Error starting server. %w", err)
		}
	case <-quit:
		log.Println("Shutting down server...")
	}

	err = s.server.Shutdown(context.Background())
	if err != nil {
		err = fmt.Errorf("Error shutting down server. %w", err)

		errStop := s.Stop()
		if errStop != nil {
			err = fmt.Errorf("Error stopping server. %w. %w", errStop, err)
		}

		return err
	}

	return nil
}

func (s *Server) Stop() error {
	return nil
}
