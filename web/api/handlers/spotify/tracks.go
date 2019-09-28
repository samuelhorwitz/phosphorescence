package spotify

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/go-chi/chi"
	"github.com/samuelhorwitz/phosphorescence/api/common"
	"net/http"
	"strings"
	"time"
)

func Tracks(w http.ResponseWriter, r *http.Request) {
	region := strings.ToLower(chi.URLParam(r, "region"))
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
