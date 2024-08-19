package task

import (
	"github.com/robfig/cron/v3"
)

var job = cron.New()

type CronTask struct {
}
