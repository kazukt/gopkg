package server

import (
	"context"
	"net/http"
	"time"
)

type Server struct {
	server          *http.Server
	shutdownTimeout time.Duration
}

func (s *Server) Run(ctx context.Context) error {
	errCh := make(chan error, 1)
	go func() {
		errCh <- s.server.ListenAndServe()
	}()

	for {
		select {
		case err := <-errCh:
			return err
		case <-ctx.Done():
			return s.shutdown()
		}
	}
}

func (s *Server) shutdown() error {
	if s.server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
		defer cancel()

		return s.server.Shutdown(ctx)
	}
	return nil
}
