package config

import "github.com/spf13/viper"

// Init load configurations from config.yml file
func Init(cfgFile string) error {
	configName := "config"
	if cfgFile != "" {
		configName = cfgFile
	}

	viper.SetConfigName(configName)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	initConfig()
	return nil
}

// initConfig laod all configurations
func initConfig() {
	loadApp()
	loadScheduler()
	loadRedis()
	loadDatabase()
}
