package tilecache

import (
	"path/filepath"
	"strconv"

	"github.com/enixlan/tiler/internal/app/core/domain"
)

func key(req *domain.TileRequest) string {
	const base = 10

	return filepath.Join(
		req.TileSetName(),
		strconv.FormatUint(req.Zoom, base),
		strconv.FormatUint(req.X, base),
		strconv.FormatUint(req.Y, base),
		"tile."+req.Format.String(),
	)
}
