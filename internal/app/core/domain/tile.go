package domain

import (
	"context"
	"fmt"
)

type Tile struct {
	Data   []byte
	Format ImageFormat
}

func (t *Tile) IsValid() bool {
	return len(t.Data) > 0 && t.Format == ImageFormatPNG
}

type TileRequest struct {
	Zoom   uint64
	X      uint64
	Y      uint64
	Format ImageFormat
}

func (r *TileRequest) TileSetName() string {
	return "default"
}

type TileRequestCtx struct {
	Request    *TileRequest
	RequestCtx context.Context

	done chan struct{}
	tile *Tile
	err  error
}

func NewTileRequestCtx(ctx context.Context, req *TileRequest) *TileRequestCtx {
	return &TileRequestCtx{
		Request:    req,
		RequestCtx: ctx,
		done:       make(chan struct{}, 1),
		tile:       nil,
		err:        nil,
	}
}

func (r *TileRequestCtx) WaitResult(ctx context.Context) (*Tile, error) {
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("select result: %w", ctx.Err())
	case <-r.done:
		return r.tile, r.err
	}
}

func (r *TileRequestCtx) IsCanceled() bool {
	select {
	case <-r.RequestCtx.Done():
		return true
	default:
		return false
	}
}

func (r *TileRequestCtx) WriteResult(tile *Tile, err error) {
	r.tile = tile
	r.err = err

	close(r.done)
}
