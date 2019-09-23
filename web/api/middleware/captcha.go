package middleware

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/didip/tollbooth"
	"github.com/samuelhorwitz/phosphorescence/api/common"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func Captcha(action string, threshold float64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			lmtErr := tollbooth.LimitByKeys(ipLimiter, []string{r.RemoteAddr})
			if lmtErr != nil {
				ipLimiter.ExecOnLimitReached(w, r)
				common.Fail(w, fmt.Errorf("Phosphorescence API IP rate limiting hit: %s", lmtErr.Message), lmtErr.StatusCode)
				return
			}
			captcha := r.URL.Query().Get("captcha")
			if captcha == "" {
				common.Fail(w, errors.New("Must include Recaptcha token"), http.StatusBadRequest)
				return
			}
			form := url.Values{}
			form.Add("secret", recaptchaSecret)
			form.Add("response", captcha)
			form.Add("remoteip", r.RemoteAddr)
			req, err := http.NewRequestWithContext(r.Context(), "POST", "https://www.google.com/recaptcha/api/siteverify", strings.NewReader(form.Encode()))
			if err != nil {
				common.Fail(w, fmt.Errorf("Could not build Recaptcha request: %s", err), http.StatusInternalServerError)
				return
			}
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			res, err := googleHTTPClient.Do(req)
			if err != nil {
				common.Fail(w, fmt.Errorf("Could not make Recaptcha request: %s", err), http.StatusInternalServerError)
				return
			}
			defer res.Body.Close()
			if res.StatusCode != http.StatusOK {
				common.Fail(w, fmt.Errorf("Recaptcha request responded with %d", res.StatusCode), http.StatusInternalServerError)
				return
			}
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				common.Fail(w, fmt.Errorf("Could not read Recaptcha response: %s", err), http.StatusInternalServerError)
				return
			}
			var recaptchaResponse struct {
				Success            bool      `json:"success"`
				Score              float64   `json:"score"`
				Action             string    `json:"action"`
				ChallengeTimestamp time.Time `json:"challenge_ts"`
				Hostname           string    `json:"hostname"`
				ErrorCodes         []string  `json:"error-codes"`
			}
			err = json.Unmarshal(body, &recaptchaResponse)
			if err != nil {
				common.Fail(w, fmt.Errorf("Could not parse Recaptcha response: %s", err), http.StatusInternalServerError)
				return
			}
			if !recaptchaResponse.Success {
				common.Fail(w, errors.New("Recaptcha not successful"), http.StatusForbidden)
				return
			}
			if recaptchaResponse.Hostname != phosphorHost {
				common.Fail(w, fmt.Errorf("Invalid Recaptcha hostname: %s", recaptchaResponse.Hostname), http.StatusForbidden)
				return
			}
			if recaptchaResponse.Action != action {
				common.Fail(w, fmt.Errorf("Invalid Recaptcha action: %s", recaptchaResponse.Action), http.StatusForbidden)
				return
			}
			if recaptchaResponse.Score < threshold {
				common.Fail(w, fmt.Errorf("Likely a bot: %f", recaptchaResponse.Score), http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
