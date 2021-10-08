package ports

import (
	"github.com/enixlan/tiler/internal/app/core/domain"
)

type IMapnikAdapter interface {
	Render(req *domain.TileRequest) (*domain.Tile, error)
	Close()
}
