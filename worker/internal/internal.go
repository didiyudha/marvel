package internal

import (
	"github.com/robfig/cron/v3"
)

const (
	config = "@every 3s"
)

type Worker interface {
	AddFunc(worker InternalWorker)
	Run()
}

type InternalWorker interface {
	Do()
}

type workerImpl struct {
	Cron *cron.Cron
}

func NewWorker() Worker {
	c := cron.New(
		cron.WithSeconds(),
		cron.WithChain(
			cron.Recover(cron.DefaultLogger),
		),
	)
	return &workerImpl{Cron: c}
}

func (w *workerImpl) AddFunc(worker InternalWorker) {
	w.Cron.AddFunc(config, worker.Do)
}

func (w *workerImpl) Run() {
	w.Cron.Run()
}