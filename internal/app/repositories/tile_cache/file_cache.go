package tilecache

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/enixlan/tiler/internal/app/core/domain"
)

type FileCache struct {
	mu  sync.RWMutex
	dir string
}

func NewFileCache(dir string) (*FileCache, error) {
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return nil, fmt.Errorf("makdir all: %w", err)
	}

	cache := &FileCache{
		dir: dir,
	}

	return cache, nil
}

func (c *FileCache) Get(_ context.Context, req *domain.TileRequest) (*domain.Tile, error) {
	const op = "FileCache.Get"

	c.mu.RLock()
	defer c.mu.RUnlock()

	path := filepath.Join(c.dir, key(req))

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrTileNotFound
		}

		return nil, fmt.Errorf("%s: read file: %w", op, err)
	}

	tile := &domain.Tile{
		Data:   data,
		Format: req.Format,
	}

	return tile, nil
}

func (c *FileCache) Set(_ context.Context, req *domain.TileRequest, tile *domain.Tile) error {
	const op = "FileCache.Set"

	c.mu.RLock()
	defer c.mu.RUnlock()

	path := filepath.Join(c.dir, key(req))

	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		return fmt.Errorf("%s: mkdir all: %w", op, err)
	}

	if err := os.WriteFile(path, tile.Data, os.ModePerm); err != nil {
		return fmt.Errorf("%s: write file: %w", op, err)
	}

	return nil
}

func (c *FileCache) Del(_ context.Context, req *domain.TileRequest) error {
	const op = "FileCache.Del"

	c.mu.RLock()
	defer c.mu.RUnlock()

	path := filepath.Join(c.dir, key(req))

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	}

	if err := os.Remove(path); err != nil {
		return fmt.Errorf("%s: remove file: %w", op, err)
	}

	return nil
}

func (c *FileCache) Reset() error {
	const op = "FileCache.Reset"

	c.mu.Lock()
	defer c.mu.Unlock()

	if _, err := os.Stat(c.dir); os.IsNotExist(err) {
		return nil
	}

	if err := os.RemoveAll(c.dir); err != nil {
		return fmt.Errorf("%s: remove all: %w", op, err)
	}

	return nil
}

func (c *FileCache) Close() error {
	return nil
}
