package queue

import (
	"context"
	"errors"
	"sync"

	"github.com/enixlan/tiler/internal/app/core/domain"
)

var ErrIsFull = errors.New("queue is full")

type Queue struct {
	ch chan *domain.TileRequestCtx
	sn sync.Once
}

func NewQueue(capacity int) *Queue {
	return &Queue{
		ch: make(chan *domain.TileRequestCtx, capacity),
		sn: sync.Once{},
	}
}

func (q *Queue) Push(ctx context.Context, req *domain.TileRequest) (*domain.TileRequestCtx, error) {
	result := domain.NewTileRequestCtx(ctx, req)

	select {
	case q.ch <- result:
		return result, nil
	default:
		return nil, ErrIsFull
	}
}

func (q *Queue) Read(cancel <-chan struct{}) (*domain.TileRequestCtx, bool) {
	select {
	case r, ok := <-q.ch:
		return r, ok
	case <-cancel:
		return nil, false
	}
}

func (q *Queue) Close() error {
	q.sn.Do(q.close)

	return nil
}

func (q *Queue) close() {
	close(q.ch)
}
