package middleware

import (
	"fmt"
	"github.com/didip/tollbooth"
	"github.com/samuelhorwitz/phosphorescence/api/common"
	"net/http"
)

func IPLimiter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lmtErr := tollbooth.LimitByKeys(ipLimiter, []string{r.RemoteAddr})
		if lmtErr != nil {
			ipLimiter.ExecOnLimitReached(w, r)
			common.Fail(w, fmt.Errorf("Phosphorescence API IP rate limiting hit: %s", lmtErr.Message), lmtErr.StatusCode)
			return
		}
		next.ServeHTTP(w, r)
	})
}
