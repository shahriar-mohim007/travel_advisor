package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"
	"travel_advisor/pkg/config"
	"travel_advisor/pkg/conn"
	"travel_advisor/pkg/log"

	travelHandler "travel_advisor/travel/delivery/http"
	travelUsecase "travel_advisor/travel/usecase"

	userHandler "travel_advisor/user/delivery/http"
	userReposiotry "travel_advisor/user/repository"
	userUsecase "travel_advisor/user/usecase"

	districtRepository "travel_advisor/districts/repository"

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
			fmt.Println("--------Database is connecting-------")
			err := conn.ConnectDefaultDB()
			if err != nil {
				log.Fatal(err)
			}
			log.Info("Database connected successfully!")

			log.Info("--------Connecting cache server--------")
			if err := conn.ConnectDefaultCache(); err != nil {
				log.Fatal(err)
			}
			log.Info("Cache server connected successfully!")

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

	db := conn.DefaultDB()
	cacher := conn.DefaultCache()

	dis := districtRepository.NewDistrictPostgreSQL(db)
	tc := travelUsecase.NewTravelUsecase(cacher, dis)
	us := userReposiotry.NewUserPostgreSQL(db)
	uc := userUsecase.NewUserUsecase(us)

	travelHandler.NewTravelHandler(r, tc)
	userHandler.NewUserHandler(r, uc)

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
