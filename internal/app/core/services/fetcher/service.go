package fetcher

import (
	"context"
	"fmt"

	"github.com/loghole/tracing/tracelog"

	"github.com/enixlan/tiler/internal/app/core/domain"
	"github.com/enixlan/tiler/internal/app/core/ports"
)

type Service struct {
	logger tracelog.Logger
	queue  ports.ITileRequestQueue
	cache  ports.ITileCache
}

func NewService(
	logger tracelog.Logger,
	queue ports.ITileRequestQueue,
	cache ports.ITileCache,
) *Service {
	return &Service{
		logger: logger,
		queue:  queue,
		cache:  cache,
	}
}

func (s *Service) FetchTile(ctx context.Context, req *domain.TileRequest) (*domain.Tile, error) {
	const op = "FetcherService.FetchTile"

	if tile, err := s.cache.Get(ctx, req); err == nil && tile.IsValid() {
		return tile, nil
	}

	requestCtx, err := s.queue.Push(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("%s: push tile req to queue: %w", op, err)
	}

	// Ожидание результатов рендеринга тайла или ошибки.
	tile, err := requestCtx.WaitResult(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: wait tile result: %w", op, err)
	}

	if err := s.cache.Set(ctx, req, tile); err != nil {
		s.logger.Errorf(ctx, "%s: store tile to cache: %v", op, err)
	}

	return tile, nil
}
