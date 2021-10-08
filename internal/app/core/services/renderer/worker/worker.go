package worker

import (
	"context"
	"sync"

	"github.com/loghole/tracing/tracelog"

	"github.com/enixlan/tiler/internal/app/core/domain"
	"github.com/enixlan/tiler/internal/app/core/ports"
)

type Worker struct {
	logger  tracelog.Logger
	queue   ports.ITileRequestQueue
	adapter ports.IMapnikAdapter

	wg      sync.WaitGroup
	closeCh chan struct{}
}

func New(logger  tracelog.Logger,queue ports.ITileRequestQueue, adapter ports.IMapnikAdapter) *Worker {
	return &Worker{
		logger:  logger,
		queue:   queue,
		adapter: adapter,
		closeCh: make(chan struct{}),
	}
}

func (w *Worker) Run() {
	w.logger.Infof(context.TODO(), "renderer worker started")

	w.wg.Add(1)
	defer w.wg.Done()

	defer w.adapter.Close()

	for {
		req, ok := w.queue.Read(w.closeCh)
		if !ok { // if queue is closed.
			break
		}

		w.render(req)
	}

	w.logger.Infof(context.TODO(), "renderer worker closed")
}

func (w *Worker) Close() {
	close(w.closeCh)
	w.wg.Wait()
}

func (w *Worker) render(req *domain.TileRequestCtx) {
	if req.IsCanceled() {
		w.logger.Debugf(
			req.RequestCtx,
			"render tile zoom=%d x=%d y=%d: canceled",
			req.Request.Zoom,
			req.Request.X,
			req.Request.Y,
		)

		req.WriteResult(nil, context.Canceled)

		return
	}

	tile, err := w.adapter.Render(req.Request)
	if err != nil {
		w.logger.Errorf(
			req.RequestCtx,
			"render tile zoom=%d x=%d y=%d: %v",
			req.Request.Zoom,
			req.Request.X,
			req.Request.Y,
			err,
		)
	} else {
		w.logger.Debugf(
			req.RequestCtx,
			"render tile zoom=%d x=%d y=%d: completed",
			req.Request.Zoom,
			req.Request.X,
			req.Request.Y,
		)
	}

	req.WriteResult(tile, err)
}
