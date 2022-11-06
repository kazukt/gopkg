package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"golang.org/x/sync/errgroup"
)

func ExampleServer_Run() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	group, groupCtx := errgroup.WithContext(ctx)

	adminMux := http.NewServeMux()
	adminMux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status":"OK"}`))
	})
	adminServer := Server{
		server: &http.Server{
			Addr:    ":10000",
			Handler: adminMux,
		},
		shutdownTimeout: 1,
	}

	appMux := http.NewServeMux()
	appMux.HandleFunc("/foo", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"Hello":"World!"}`))
	})
	appServer := Server{
		server: &http.Server{
			Addr:    ":8080",
			Handler: appMux,
		},
		shutdownTimeout: 1,
	}

	group.Go(func() error {
		return adminServer.Run(groupCtx)
	})

	group.Go(func() error {
		return appServer.Run(groupCtx)
	})

	if err := group.Wait(); err != nil {
		fmt.Println("unexpected error")
	}
}

func TestServer_Run(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	srv := Server{
		server: &http.Server{
			Addr:    ":8080",
			Handler: http.DefaultServeMux,
		},
		shutdownTimeout: 1 * time.Second,
	}
	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.Run(ctx)
	}()

	if err := <-errCh; err != nil {
		t.Error(err)
	}
}
