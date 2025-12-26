package cmd

import (
	"os"
	"travel_advisor/pkg/config"
	"travel_advisor/pkg/log"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	cfgFile                 string
	verbose, prettyPrintLog bool
	rootCmd                 = &cobra.Command{
		Use:   "travel_advisor",
		Short: "travel_advisor",
		Long:  `travel_advisor`,
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(InitConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "config.yml", "config file")
	rootCmd.PersistentFlags().BoolVarP(&prettyPrintLog, "pretty", "p", false, "pretty print verbose/log")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
}

func InitConfig() {
	log.Info("Loading configurations")
	if err := config.Init(cfgFile); err != nil {
		log.Warn("Failed to load configuration")
		log.Fatal(err)
	}
	log.Info("Configurations loaded successfully!")

	// Log as JSON instead of the default ASCII formatter.
	log.SetLogFormatter(&logrus.JSONFormatter{
		PrettyPrint: prettyPrintLog,
	})

	log.SetLogLevel(logrus.TraceLevel)
	if verbose {
		log.SetLogLevel(logrus.TraceLevel)
	}
}
