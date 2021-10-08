package tilecache

import (
	"context"
	"fmt"

	"github.com/loghole/tracing/tracelog"

	"github.com/enixlan/tiler/internal/app/core/domain"
	"github.com/enixlan/tiler/internal/app/core/ports"
)

type fillEntry struct {
	layer int
	req   *domain.TileRequest
	tile  *domain.Tile
}

type Wrapper struct {
	logger tracelog.Logger
	layers []ports.ITileCache

	fillCh chan *fillEntry
}

func NewWrapper(logger tracelog.Logger, layers ...ports.ITileCache) *Wrapper {
	const capacity = 1000

	wrapper := &Wrapper{
		logger: logger,
		layers: layers,
		fillCh: make(chan *fillEntry, capacity),
	}

	go wrapper.fillWorker()

	return wrapper
}

func (w *Wrapper) Get(ctx context.Context, req *domain.TileRequest) (*domain.Tile, error) {
	var errors []error

	for idx, layer := range w.layers {
		tile, err := layer.Get(ctx, req)
		if err != nil {
			errors = append(errors, fmt.Errorf("layer %d: %w", idx, err))

			continue
		}

		w.fillUpperLayers(idx, req, tile)

		return tile, nil
	}

	return nil, fmt.Errorf("%w: [%v]", ErrTileNotFound, errors)
}

func (w *Wrapper) Set(ctx context.Context, req *domain.TileRequest, tile *domain.Tile) error {
	var errors []error

	for idx, layer := range w.layers {
		if err := layer.Set(ctx, req, tile); err != nil {
			errors = append(errors, fmt.Errorf("layer %d: %w", idx, err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("%w: [%v]", ErrWriteTileFailed, errors)
	}

	return nil
}

func (w *Wrapper) Reset() error {
	return nil
}

func (w *Wrapper) Close() error {
	close(w.fillCh)

	var errors []error

	for idx, layer := range w.layers {
		if err := layer.Close(); err != nil {
			errors = append(errors, fmt.Errorf("layer %d: %w", idx, err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("%w: [%v]", ErrCloseFailed, errors)
	}

	return nil
}

func (w *Wrapper) fillWorker() {
	for entry := range w.fillCh {
		for i := entry.layer - 1; i >= 0; i-- {
			if err := w.layers[i].Set(context.TODO(), entry.req, entry.tile); err != nil {
				w.logger.Errorf(context.TODO(), "fill layer %d : %v", i, err)
			}
		}
	}
}

func (w *Wrapper) fillUpperLayers(layer int, req *domain.TileRequest, tile *domain.Tile) {
	select {
	case w.fillCh <- &fillEntry{layer: layer, req: req, tile: tile}:
	default:
		w.logger.Error(context.TODO(), "fill upper layers skipped: queue is full")
	}
}
