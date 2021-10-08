package mapnik

import (
	"errors"
	"fmt"

	mapnik "github.com/gadavy/go-mapnik"

	"github.com/enixlan/tiler/internal/app/core/domain"
	"github.com/enixlan/tiler/pkg/googleprojection"
)

var ErrInvalidImageFormat = errors.New("invalid image format")

const (
	maxZoom    = 30
	tileSize   = 256
	bufferSize = 128
)

func Init() error {
	if err := mapnik.RegisterFonts(mapnik.ConfigFonts()); err != nil {
		return fmt.Errorf("init mapnik fonts: %w", err)
	}

	if err := mapnik.RegisterDatasources(mapnik.ConfigPlugins()); err != nil {
		return fmt.Errorf("init mapnik data sources: %w", err)
	}

	return nil
}

func Version() string {
	return mapnik.Version()
}

type Mapnik struct {
	m  *mapnik.Map
	p  *mapnik.Projection
	gp *googleprojection.Mercator
}

func NewMapnik(stylesheet string) (*Mapnik, error) {
	m := mapnik.NewMap(tileSize, tileSize)
	m.SetBufferSize(bufferSize)

	if err := m.Load(stylesheet); err != nil {
		return nil, fmt.Errorf("load stylesheet: %w", err)
	}

	adapter := &Mapnik{
		m:  m,
		p:  m.Projection(),
		gp: googleprojection.New(tileSize, maxZoom),
	}

	return adapter, nil
}

func (m *Mapnik) Render(req *domain.TileRequest) (*domain.Tile, error) {
	var (
		p0x = float64(req.X) * float64(tileSize)
		p0y = (float64(req.Y) + 1) * float64(tileSize)
		p1x = (float64(req.X) + 1) * float64(tileSize)
		p1y = float64(req.Y) * float64(tileSize)
	)

	// Convert to LatLong(EPSG:4326)
	l0x, l0y, err := m.gp.PixelToLatLong(p0x, p0y, req.Zoom)
	if err != nil {
		return nil, fmt.Errorf("p0x p0y to lat long: %w", err)
	}

	l1x, l1y, err := m.gp.PixelToLatLong(p1x, p1y, req.Zoom)
	if err != nil {
		return nil, fmt.Errorf("p1x p1y to lat long: %w", err)
	}

	// Get render option (image fmt, scale, etc.)
	opts, err := m.renderOpts(req.Format)
	if err != nil {
		return nil, fmt.Errorf("render opts: %w", err)
	}

	// Convert to map projection (e.g. mercartor co-ords EPSG:3857)
	c0x, c0y := m.p.Forward(l0x, l0y)
	c1x, c1y := m.p.Forward(l1x, l1y)

	// Zoom and render tile.
	m.m.Resize(tileSize, tileSize)
	m.m.Zoom(c0x, c0y, c1x, c1y)

	data, err := m.m.Render(opts)
	if err != nil {
		return nil, fmt.Errorf("render to %s: %w", opts.Format, err)
	}

	tile := &domain.Tile{
		Data:   data,
		Format: req.Format,
	}

	return tile, err
}

func (m *Mapnik) Close() {
	m.p.Free()
	m.m.Free()
}

func (m *Mapnik) renderOpts(format domain.ImageFormat) (opts mapnik.RenderOpts, err error) {
	switch format {
	case domain.ImageFormatPNG:
		opts.Format = mapnik.Png256
	default:
		return opts, ErrInvalidImageFormat
	}

	return opts, nil
}
