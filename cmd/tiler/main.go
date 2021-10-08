package main

import (
	"log"

	"github.com/loghole/tron"

	"github.com/enixlan/tiler/config"
	"github.com/enixlan/tiler/internal/app/adapters/mapnik"
	"github.com/enixlan/tiler/internal/app/adapters/queue"
	"github.com/enixlan/tiler/internal/app/api/tile"
	"github.com/enixlan/tiler/internal/app/core/services/fetcher"
	"github.com/enixlan/tiler/internal/app/core/services/renderer"
	tilecache "github.com/enixlan/tiler/internal/app/repositories/tile_cache"
)

func main() {
	app, err := tron.New(tron.AddLogCaller(), tron.WithRealtimeConfig())
	if err != nil {
		log.Fatalf("can't create app: %s", err)
	}

	defer app.Close()

	logger := app.TraceLogger()

	// Init mapnik stylesheet.
	stylesheet, err := config.MapnikStylesheet()
	if err != nil {
		app.Logger().Fatalf("init mapnik stylesheet: %s", err)
	}

	// Init adapters.
	if err := mapnik.Init(); err != nil {
		app.Logger().Fatalf("init mapnik fonts: %s", err)
	}

	app.Logger().Infof("mapnik version %s", mapnik.Version())

	queueCap, err := config.GetRendererQueueCapacity()
	if err != nil {
		app.Logger().Fatalf("get renderer queue cap: %s", err)
	}

	requestQueue := queue.NewQueue(queueCap)

	// Init caches.
	memoryCacheMaxSize, err := config.GetMemoryCacheMaxSize()
	if err != nil {
		app.Logger().Fatalf("get tile memory cache max size: %s", err)
	}

	memoryCacheLifeWindow, err := config.GetMemoryCacheLifeWindow()
	if err != nil {
		app.Logger().Fatalf("get tile memory cache life window: %s", err)
	}

	memoryCache, err := tilecache.NewMemoryCache(memoryCacheMaxSize, memoryCacheLifeWindow)
	if err != nil {
		app.Logger().Fatalf("init tile cache: %v", err)
	}

	fileCacheDir, err := config.GetFileCacheDir()
	if err != nil {
		app.Logger().Fatalf("get tile file cache dir: %s", err)
	}

	fileCache, err := tilecache.NewFileCache(fileCacheDir)
	if err != nil {
		app.Logger().Fatalf("init tile cache: %v", err)
	}

	wrappedCache := tilecache.NewWrapper(logger, memoryCache, fileCache)

	workersCount, err := config.GetRendererWorkersCount()
	if err != nil {
		app.Logger().Fatalf("get renderer workers count: %s", err)
	}

	// Init services.
	var (
		fetcherService  = fetcher.NewService(logger, requestQueue, wrappedCache)
		rendererService = renderer.NewService(logger, requestQueue, stylesheet)
	)

	config.InitWorkersCountSetter(rendererService.SetWorkersCount)

	// Init handlers.
	tileAPI := tile.NewImplementation(logger, fetcherService)

	app.Router().Get("/map/", tileAPI.MapHandler)
	app.Router().Get("/tile/default/{zoom}/{x}/{y}.png", tileAPI.TileHandler)

	if err := rendererService.Run(workersCount); err != nil {
		app.Logger().Fatalf("init renderer service: %v", err)
	}

	// Run application.
	if err := app.Run(); err != nil {
		app.Logger().Fatalf("can't run app: %v", err)
	}

	// Close all.
	if err := requestQueue.Close(); err != nil {
		app.Logger().Errorf("close request queue: %v", err)
	}

	if err := rendererService.Close(); err != nil {
		app.Logger().Errorf("close renderer service: %v", err)
	}

	if err := wrappedCache.Close(); err != nil {
		app.Logger().Errorf("close wrapped cache: %v", err)
	}
}
