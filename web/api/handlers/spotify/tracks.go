package spotify

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/go-chi/chi"
	"github.com/samuelhorwitz/phosphorescence/api/common"
	"github.com/samuelhorwitz/phosphorescence/api/middleware"
	"github.com/samuelhorwitz/phosphorescence/api/session"
)

func TracksUnauthenticated(w http.ResponseWriter, r *http.Request) {
	tracks(w, r, strings.ToLower(chi.URLParam(r, "region")))
}

func Tracks(w http.ResponseWriter, r *http.Request) {
	sess, ok := r.Context().Value(middleware.AuthenticatedSessionContextKey).(*session.Session)
	if !ok {
		common.Fail(w, errors.New("No session on request context"), http.StatusUnauthorized)
		return
	}
	tracks(w, r, sess.SpotifyCountry)
}

func tracks(w http.ResponseWriter, r *http.Request, region string) {
	s3Req, _ := s3Service.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String("phosphorescence-tracks"),
		Key:    aws.String(fmt.Sprintf("tracks.%s.json", region)),
	})
	var err error
	tracksURL, err := s3Req.Presign(2 * time.Minute)
	if err != nil {
		common.Fail(w, fmt.Errorf("Could not presign Spaces request: %s", err), http.StatusInternalServerError)
		return
	}
	common.JSON(w, struct {
		TracksURL string `json:"tracksUrl"`
	}{
		TracksURL: tracksURL,
	})
}
