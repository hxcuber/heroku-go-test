package httpserv

import (
	"context"
	"fmt"
	"github.com/hxcuber/friends-management/api/pkg/env"
	"github.com/pkg/errors"
	"net/http"
	"os"
	"time"

	pkgerrors "github.com/pkg/errors"
)

// server represents an HTTP server
type server struct {
	srv           *http.Server
	shutdownGrace time.Duration
}

// NewServer initializes and returns an instance of HTTP Server
func NewServer(handler http.Handler) *server {
	s := server{
		srv: &http.Server{
			Addr:         ":" + env.GetAndValidateF("PORT"),
			Handler:      handler,
			ReadTimeout:  time.Minute,
			WriteTimeout: time.Minute,
		},
		shutdownGrace: 10 * time.Second,
	}

	return &s
}

// Start starts the HTTP server
func (s *server) Start(ctx context.Context) error {
	// Make a channel to listen for errors coming from the listener. Use a
	// buffered channel so the goroutine can exit if we don't collect this error.
	startupError := make(chan error, 1)

	// Start the service listening for requests.
	go func() {
		startupError <- s.srv.ListenAndServe()
	}()

	fmt.Fprintf(os.Stdout, "Running locally on: http://localhost%s/\n", s.srv.Addr)

	// Blocking main and waiting for shutdown.
	select {
	case err := <-startupError:
		if !errors.Is(err, http.ErrServerClosed) { // ListenAndServe will always return a non-nil error
			return pkgerrors.WithStack(fmt.Errorf("startup failed: %w", err))
		}
		return nil
	case <-ctx.Done():
		return s.Stop()
	}
}

// Stop stops the HTTP server
func (s *server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownGrace)
	defer cancel()

	if err := s.srv.Shutdown(ctx); err != nil {
		if err = s.srv.Close(); err != nil {
			return pkgerrors.WithStack(fmt.Errorf("force shutdown failed: %w", err))
		}
	}

	return nil
}
