package middleware

import (
	"context"
	"math"
	"net/http"
	"strconv"
	"time"
)

const PageCountContextKey = contextKey("pageCount")
const PageCursorContextKey = contextKey("pageCursor")
const maxPageSize = 10

func Paginate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count := uint64(maxPageSize)
		if countStr := r.URL.Query().Get("count"); countStr != "" {
			if tentativeCount, err := strconv.ParseUint(countStr, 10, 64); err == nil {
				count = uint64(math.Max(1, math.Min(float64(tentativeCount), maxPageSize)))
			}
		}
		var from time.Time
		if fromStr := r.URL.Query().Get("from"); fromStr != "" {
			if tentativeFrom, err := time.Parse(time.RFC3339, fromStr); err == nil {
				from = tentativeFrom
			}
		}
		ctx := r.Context()
		ctx = context.WithValue(ctx, PageCountContextKey, count)
		ctx = context.WithValue(ctx, PageCursorContextKey, from)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
