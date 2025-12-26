package config

import (
	"time"

	"github.com/spf13/viper"
)

type HttpApplication struct {
	HTTPPort int
	Verbose  bool

	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration

	PaginationLimit int
}

type AppConfig struct {
	JwtSecret string
}

var http_app HttpApplication
var appConfig AppConfig

func HttpApp() HttpApplication {
	return http_app
}

func App() AppConfig {
	return appConfig
}

func loadApp() {
	appConfig = AppConfig{
		JwtSecret: viper.GetString("app.jwt_secret"),
	}

	http_app = HttpApplication{
		HTTPPort:        viper.GetInt("http_app.http_port"),
		Verbose:         viper.GetBool("http_app.verbose"),
		ReadTimeout:     viper.GetDuration("http_app.read_timeout") * time.Second,
		WriteTimeout:    viper.GetDuration("http_app.write_timeout") * time.Second,
		IdleTimeout:     viper.GetDuration("http_app.idle_timeout") * time.Second,
		PaginationLimit: viper.GetInt("http_app.pagination_limit"),
	}
}
