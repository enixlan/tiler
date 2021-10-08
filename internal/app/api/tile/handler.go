package tile

import (
	"github.com/loghole/tracing/tracelog"

	"github.com/enixlan/tiler/internal/app/core/ports"
)

type Implementation struct {
	logger  tracelog.Logger
	service ports.IFetcherService
}

func NewImplementation(logger  tracelog.Logger, service ports.IFetcherService) *Implementation {
	return &Implementation{
		logger: logger,
		service: service,
	}
}
