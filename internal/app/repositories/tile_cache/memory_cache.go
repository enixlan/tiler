package tilecache

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/allegro/bigcache/v3"

	"github.com/enixlan/tiler/internal/app/core/domain"
)

const (
	_MB = 1 << (10 * 2)

	_shards      = 1024
	_cleanWindow = time.Minute
	_maxTileSize = 30 * 1024 // 30kB должно хватать, обычно тайл весит ~ 10kB

	_defaultMaxSizeMB = 1024
	_minimalMaxSizeMB = 10
)

type MemoryCache struct {
	cache  *bigcache.BigCache
	config bigcache.Config
}

func NewMemoryCache(maxSizeMB int, lifeWindow time.Duration) (*MemoryCache, error) {
	cache := &MemoryCache{}

	if maxSizeMB < _minimalMaxSizeMB {
		maxSizeMB = _defaultMaxSizeMB
	}

	cache.config = bigcache.Config{
		Shards:             _shards,
		LifeWindow:         lifeWindow,
		CleanWindow:        _cleanWindow,
		MaxEntriesInWindow: maxSizeMB * _MB / _maxTileSize,
		MaxEntrySize:       _maxTileSize,
		HardMaxCacheSize:   maxSizeMB,
	}

	var err error

	cache.cache, err = bigcache.NewBigCache(cache.config)
	if err != nil {
		return nil, fmt.Errorf("init chache: %w", err)
	}

	return cache, nil
}

func (c *MemoryCache) Get(_ context.Context, req *domain.TileRequest) (*domain.Tile, error) {
	const op = "MemoryCache.Get"

	data, err := c.cache.Get(key(req))
	if err != nil {
		if errors.Is(err, bigcache.ErrEntryNotFound) {
			return nil, ErrTileNotFound
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	tile := &domain.Tile{
		Data:   data,
		Format: req.Format,
	}

	return tile, nil
}

func (c *MemoryCache) Set(_ context.Context, req *domain.TileRequest, tile *domain.Tile) error {
	const op = "MemoryCache.Set"

	if err := c.cache.Set(key(req), tile.Data); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (c *MemoryCache) Del(_ context.Context, req *domain.TileRequest) error {
	const op = "MemoryCache.Del"

	if err := c.cache.Delete(key(req)); err != nil {
		if errors.Is(err, bigcache.ErrEntryNotFound) {
			return nil
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (c *MemoryCache) Reset() error {
	const op = "MemoryCache.Reset"

	if err := c.cache.Reset(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (c *MemoryCache) Close() error {
	const op = "MemoryCache.Close"

	if err := c.cache.Close(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
