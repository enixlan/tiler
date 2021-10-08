package renderer

import (
	"context"
	"fmt"
	"sync"

	"github.com/loghole/tracing/tracelog"

	"github.com/enixlan/tiler/internal/app/adapters/mapnik"
	"github.com/enixlan/tiler/internal/app/core/ports"
	"github.com/enixlan/tiler/internal/app/core/services/renderer/worker"
)

type Service struct {
	logger     tracelog.Logger
	queue      ports.ITileRequestQueue
	stylesheet string

	mu         sync.Mutex
	workerPool []*worker.Worker
}

func NewService(logger tracelog.Logger, queue ports.ITileRequestQueue, stylesheet string) *Service {
	return &Service{
		logger:     logger,
		queue:      queue,
		stylesheet: stylesheet,
	}
}

func (s *Service) Run(count int) error {
	const op = "RendererService.Run"

	if err := s.SetWorkersCount(count); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Service) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, wr := range s.workerPool {
		wr.Close()
	}

	return nil
}

func (s *Service) SetWorkersCount(count int) (err error) {
	const op = "RendererService.SetWorkersCount"

	s.mu.Lock()
	defer s.mu.Unlock()

	if count < 0 {
		count = 0
	}

	switch {
	case len(s.workerPool) < count:
		err = s.addWorkers(count-len(s.workerPool), s.stylesheet)
	case len(s.workerPool) > count:
		err = s.delWorkers(len(s.workerPool) - count)
	}

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Service) addWorkers(count int, stylesheet string) error {
	for i := 0; i < count; i++ {
		s.logger.Info(context.TODO(), "init new renderer worker")

		adapter, err := mapnik.NewMapnik(stylesheet)
		if err != nil {
			return fmt.Errorf("new mapnik: %w", err)
		}

		wr := worker.New(s.logger, s.queue, adapter)
		go wr.Run()

		s.workerPool = append(s.workerPool, wr)
	}

	return nil
}

func (s *Service) delWorkers(count int) error {
	for i := 0; i < count; i++ {
		var wr *worker.Worker

		wr, s.workerPool = s.workerPool[0], s.workerPool[1:]

		wr.Close()
	}

	return nil
}
