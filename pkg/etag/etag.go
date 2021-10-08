package etag

import (
	"strconv"

	"github.com/cespare/xxhash/v2"
)

func CalculateETag(b []byte) string {
	const base = 32

	return strconv.FormatInt(int64(len(b)), base) + "-" + strconv.FormatUint(xxhash.Sum64(b), base)
}
