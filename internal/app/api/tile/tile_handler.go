package tile

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"

	"github.com/enixlan/tiler/internal/app/core/domain"
	"github.com/enixlan/tiler/pkg/etag"
)

const (
	ifNoneMatchHeader = "If-None-Match"
	etagHeader        = "ETag"
)

func (i *Implementation) TileHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req, err := i.parseTileRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	result, err := i.service.FetchTile(ctx, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	eTag := etag.CalculateETag(result.Data)
	if eTag == r.Header.Get(ifNoneMatchHeader) {
		w.WriteHeader(http.StatusNotModified)

		return
	}

	// ETag что бы не гонять лишние данные если они остались прежними.
	w.Header().Add(etagHeader, eTag)

	w.Header().Add("Cache-Control", "max-age=86400")
	w.Header().Add("Cache-Control", "stale-while-revalidate=172800")
	w.Header().Add("Cache-Control", "stale-if-error=172800")

	// Должен игнорироваться т.к. есть Cache-Control, но не всегда.
	w.Header().Add("Expires", time.Now().Add(time.Hour*24).Format(time.RFC1123))

	if ct, ok := result.Format.ContentType(); ok {
		w.Header().Set("Content-Type", ct)
	} else {
		w.Header().Set("Content-Type", http.DetectContentType(result.Data))
	}

	if _, err := w.Write(result.Data); err != nil {
		i.logger.Errorf(ctx, "write tile data: %v", err)

		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (i *Implementation) parseTileRequest(r *http.Request) (*domain.TileRequest, error) {
	const op = "TileImplementation.parseTileRequest"

	const (
		base    = 10
		bitSize = 64
	)

	zoom, err := strconv.ParseUint(chi.URLParam(r, "zoom"), base, bitSize)
	if err != nil {
		return nil, fmt.Errorf("%s: parse zoom: %w", op, err)
	}

	x, err := strconv.ParseUint(chi.URLParam(r, "x"), base, bitSize)
	if err != nil {
		return nil, fmt.Errorf("%s: parse X: %w", op, err)
	}

	y, err := strconv.ParseUint(chi.URLParam(r, "y"), base, bitSize)
	if err != nil {
		return nil, fmt.Errorf("%s: parse Y: %w", op, err)
	}

	req := &domain.TileRequest{
		Zoom:   zoom,
		X:      x,
		Y:      y,
		Format: domain.ImageFormatPNG,
	}

	return req, nil
}
