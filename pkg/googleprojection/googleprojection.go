// projection
//
// Based on:
//  - https://en.wikipedia.org/wiki/Web_Mercator_projection
//  - https://mange.github.io/googleprojection-rs/src/googleprojection/src/lib.rs.html#1-230

package googleprojection

import (
	"errors"
	"fmt"
	"math"
)

const (
	degToRad = math.Pi / 180
	radToDeg = 180 / math.Pi
)

var ErrZoomExceeded = errors.New("zoom exceeded")

type Mercator struct {
	Bc []float64
	Cc []float64
	zc [][2]float64
	Ac []float64

	maxZoom uint64
}

// nolint:gomnd // коэффициенты из переписанных формул.
func New(size, maxZoom uint64) *Mercator {
	projection := &Mercator{maxZoom: maxZoom}

	c := float64(size)

	for d := uint64(0); d < maxZoom; d++ {
		e := c / 2
		projection.Bc = append(projection.Bc, c/360.0)
		projection.Cc = append(projection.Cc, c/(2*math.Pi))
		projection.zc = append(projection.zc, [2]float64{e, e})
		projection.Ac = append(projection.Ac, c)
		c *= 2
	}

	return projection
}

func (p *Mercator) PixelToLatLong(x, y float64, zoom uint64) (lat, long float64, err error) {
	if zoom > p.maxZoom {
		return 0, 0, fmt.Errorf("%w: max zoom = '%d'", ErrZoomExceeded, p.maxZoom)
	}

	var (
		e = p.zc[zoom]
		g = (y - e[1]) / -p.Cc[zoom]
	)

	lat = (x - e[0]) / p.Bc[zoom]
	long = radToDeg * (2*math.Atan(math.Exp(g)) - 0.5*math.Pi)

	return lat, long, nil
}

// nolint:gomnd // коэффициенты из переписанных формул.
func (p *Mercator) LatLongToPixel(x, y float64, zoom uint64) (lat, long float64, err error) {
	if zoom > p.maxZoom {
		return 0, 0, fmt.Errorf("%w: max zoom = '%d'", ErrZoomExceeded, p.maxZoom)
	}

	var (
		d = p.zc[zoom]
		f = minmax(math.Sin(y*degToRad), -0.9999, 0.9999)
	)

	lat = math.Trunc((d[0] + x*p.Bc[zoom]) + 0.5)
	long = math.Trunc((d[1] + 0.5*math.Log((1+f)/(1-f))*-p.Cc[zoom]) + 0.5)

	return lat, long, nil
}

func minmax(a, b, c float64) float64 {
	a = math.Max(a, b)
	a = math.Min(a, c)

	return a
}
