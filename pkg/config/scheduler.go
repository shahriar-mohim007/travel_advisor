package config

import (
	"github.com/spf13/viper"
)

type SchedulerCfg struct {
	CronExpr string `json:"cron_expr"`
}

var scheduler SchedulerCfg

// Scheduler contains Scheduler configurations
func Scheduler() SchedulerCfg {
	return scheduler
}

func loadScheduler() {
	scheduler = SchedulerCfg{
		CronExpr: viper.GetString("scheduler.cron_expr"),
	}
}
