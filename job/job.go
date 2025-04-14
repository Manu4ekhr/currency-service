package job

import (
	"time"

	"github.com/go-co-op/gocron/v2"
)

type Task func()

type Scheduler interface {
	Register(interval time.Duration, task Task) error
	Start()
	Shutdown()
}

type scheduler struct {
	scheduler gocron.Scheduler
}

func NewScheduler() (Scheduler, error) {
	sched, err := gocron.NewScheduler()
	if err != nil {
		return nil, err
	}

	return &scheduler{
		scheduler: sched,
	}, nil
}

func (s *scheduler) Register(interval time.Duration, task Task) error {
	_, err := s.scheduler.NewJob(
		gocron.DurationJob(interval),
		gocron.NewTask(task),
	)
	return err
}

func (s *scheduler) Start() {
	s.scheduler.Start()
}

func (s *scheduler) Shutdown() {
	s.scheduler.Shutdown()
}
