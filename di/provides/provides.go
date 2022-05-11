package provides

import (
	"stockcontent-monitor-demo-back/core/config"
	"time"
)

func ProvidesTimeout(cfg config.Config) time.Duration {
	return cfg.UseCaseTimeout()
}
