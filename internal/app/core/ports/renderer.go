package ports

import (
	"context"

	"github.com/enixlan/tiler/internal/app/core/domain"
)

type ITileRequestQueue interface {
	Push(ctx context.Context, req *domain.TileRequest) (*domain.TileRequestCtx, error)
	Read(cancel <-chan struct{}) (*domain.TileRequestCtx, bool)
}

type IFetcherService interface {
	FetchTile(ctx context.Context, req *domain.TileRequest) (*domain.Tile, error)
}

type ITileCache interface {
	Get(ctx context.Context, req *domain.TileRequest) (*domain.Tile, error)
	Set(ctx context.Context, req *domain.TileRequest, tile *domain.Tile) error
	Reset() error
	Close() error
}
