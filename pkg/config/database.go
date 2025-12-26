package config

import (
	"time"

	"github.com/spf13/viper"
)

// Database represents database configuration
type Database struct {
	Host     string              `json:"host"`
	Port     int                 `json:"port"`
	Name     string              `json:"name"`
	Username string              `json:"username"`
	Password string              `json:"password"`
	Options  map[string][]string `json:"options"`

	MaxIdleConn         int           `json:"max_idle_connection"`
	MaxOpenConn         int           `json:"max_open_connection"`
	MaxConnLifetime     time.Duration `json:"max_connection_life"`
	MaxIdleConnLifetime time.Duration `json:"max_idle_connection_life"`
	PingInterval        time.Duration `json:"ping_interval"`
}

var db Database

// DB contains database configuration
func DB() Database {
	return db
}

func loadDatabase() {
	db = Database{
		Host:     viper.GetString("postgres.host"),
		Port:     viper.GetInt("postgres.port"),
		Name:     viper.GetString("postgres.name"),
		Username: viper.GetString("postgres.username"),
		Password: viper.GetString("postgres.password"),
		Options:  viper.GetStringMapStringSlice("postgres.options"),

		MaxIdleConn:         viper.GetInt("postgres.max_idle_connection"),
		MaxOpenConn:         viper.GetInt("postgres.max_open_connection"),
		MaxConnLifetime:     viper.GetDuration("postgres.max_connection_lifetime") * time.Second,
		MaxIdleConnLifetime: viper.GetDuration("postgres.max_idle_connection_lifetime") * time.Second,

		PingInterval: viper.GetDuration("postgres.ping_interval") * time.Second,
	}
}
