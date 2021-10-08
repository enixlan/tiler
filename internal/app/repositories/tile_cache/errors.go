package tilecache

import (
	"errors"
)

var (
	ErrTileNotFound    = errors.New("tile not found in cache")
	ErrWriteTileFailed = errors.New("write tile failed")
	ErrCloseFailed     = errors.New("close failed")
)
