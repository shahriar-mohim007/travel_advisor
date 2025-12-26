package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
	"travel_advisor/dependencies"
	"travel_advisor/domain"
	"travel_advisor/helpers"
	"travel_advisor/pkg/config"
	"travel_advisor/pkg/conn"
	"travel_advisor/pkg/log"

	"github.com/robfig/cron/v3"

	"github.com/spf13/cobra"
)

var (
	schedulerCmd = &cobra.Command{
		Use:   "scheduler",
		Short: "Scheduler spawn worker process for long running background jobs",
		Long:  `Scheduler spawn worker process for long running background jobs`,
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
		RunE: scheduler,
	}
)

func init() {
	cobra.OnInitialize(InitConfig)
	rootCmd.AddCommand(schedulerCmd)
}
func scheduler(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	log.Info("Starting scheduler...")

	_, cancel := context.WithCancel(ctx)
	defer cancel()
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	cfg := config.Scheduler()
	repositories := dependencies.InjectRepositories()

	go func(ctx context.Context) {
		if err := ScheduleDistrictCacheRefresh(ctx, cfg, repositories); err != nil {
			log.Warn("failed to schedule check presigned url status cron:", err)
		}
	}(ctx)

	// Wait for the shutdown signal
	<-sigCh
	log.Warn("Shutdown signal received")

	// Optional: give background workers time to clean up
	_, shutdownCancel := context.WithTimeout(ctx, 5*time.Second)
	defer shutdownCancel()

	log.Info("-----Shutting down scheduler----")

	return nil

}

func ScheduleDistrictCacheRefresh(ctx context.Context, cfg config.SchedulerCfg, repositories dependencies.RepositoryInterfaces) error {
	s := cron.New()
	_, err := s.AddFunc(cfg.CronExpr, func() {
		districts, err := repositories.Districts.List(ctx, &domain.DistrictCriteria{})
		if err != nil {
			log.Warn("failed to get all districts")
		}
		conn.InitClient()
		client := conn.GetHTTClient()
		wg := sync.WaitGroup{}
		wg.Add(len(districts))
		for _, d := range districts {
			d := d

			go func() {
				defer wg.Done()

				log.Info("Starting district " + d.Name)

				temp, err := helpers.FetchAvgTempAt2PM(ctx, client, d.Lat, d.Long, nil)
				if err != nil {
					log.Println(err)
					log.Warn("temp fetch failed ", d.Name)
					return
				}

				pm25, err := helpers.FetchAvgPM25(ctx, client, d.Lat, d.Long, nil)
				if err != nil {
					log.Println(err)
					log.Warn("air quality fetch failed", d.Name)
					return
				}
				districtCache := domain.DistrictCache{
					Name:       d.Name,
					AvgTemp2PM: temp,
					AvgPM25:    pm25,
				}
				bytes, err := json.Marshal(districtCache)
				if err != nil {
					log.Warn("marshal failed", err)
					return
				}

				err = repositories.Cacher.Set(ctx, d.Name, bytes, time.Hour*24)
				if err != nil {
					log.Warn("failed to set cache", d.Name, err)
				}

			}()
		}

		wg.Wait()
		log.Info("district weather cache refreshed")
	})

	if err != nil {
		log.Error(ctx, "Error adding cron job: %v", err)
		return err
	}
	s.Start()

	return nil
}
