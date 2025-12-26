package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"
	"travel_advisor/pkg/config"
	"travel_advisor/pkg/log"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/spf13/cobra"
)

var (
	serveCmd = &cobra.Command{
		Use:   "serve",
		Short: "Serve run Rest server on defined port on env",
		Long:  `Serve run Rest server on defined port on env`,
		PreRun: func(cmd *cobra.Command, args []string) {

		},
		Run: serve,
	}
)

func init() {
	rootCmd.AddCommand(serveCmd)
}

func serve(cmd *cobra.Command, args []string) {
	// Initialize stop channel for graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, os.Kill)

	// build and run http server
	httpCfg := config.HttpApp()
	httpSrv := buildHTTP(cmd, args, httpCfg)
	go func() {
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("HTTP server failed: ", err)
		}
	}()

	<-stop
	log.Println("Shutting down servers...")
	// Shutdown HTTP server
	httpCtx, httpCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer httpCancel()
	if err := httpSrv.Shutdown(httpCtx); err != nil && err != http.ErrServerClosed {
		log.Fatal("Failed to shutdown HTTP server: ", err)
	}
	log.Println("Server shutdown successful!")
}

func buildHTTP(cmd *cobra.Command, args []string, httpCfg config.HttpApplication) *http.Server {
	r := chi.NewRouter()
	// middlewares
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.RealIP)

	httpPort := fmt.Sprintf(":%d", httpCfg.HTTPPort)
	log.Println("HTTP Listening on port", httpPort)

	return &http.Server{
		Addr:              httpPort,
		Handler:           r,
		ReadHeaderTimeout: httpCfg.ReadTimeout * time.Second,
		WriteTimeout:      httpCfg.WriteTimeout * time.Second,
		IdleTimeout:       httpCfg.IdleTimeout * time.Second,
	}
}
