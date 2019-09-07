package phosphor

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/samuelhorwitz/phosphorescence/api/common"
	"github.com/samuelhorwitz/phosphorescence/api/models"
)

func Search(w http.ResponseWriter, r *http.Request) {
	query := chi.URLParam(r, "query")
	if query == "" {
		common.Fail(w, errors.New("Must include search query"), http.StatusBadRequest)
		return
	}
	if len([]byte(query)) > 256 {
		common.Fail(w, errors.New("Search query too long"), http.StatusBadRequest)
		return
	}
	results, err := models.Query(query)
	if err != nil {
		common.Fail(w, fmt.Errorf("Could not execute search: %s", err), http.StatusInternalServerError)
		return
	}
	common.JSON(w, map[string]interface{}{"results": results})
}

func SearchTag(w http.ResponseWriter, r *http.Request) {
	tag := chi.URLParam(r, "tag")
	if tag == "" {
		common.Fail(w, errors.New("Must include search tag"), http.StatusBadRequest)
		return
	}
	results, err := models.QueryTag(tag)
	if err != nil {
		common.Fail(w, fmt.Errorf("Could not execute tag search: %s", err), http.StatusInternalServerError)
		return
	}
	common.JSON(w, map[string]interface{}{"results": results})
}

func RecommendedQuery(w http.ResponseWriter, r *http.Request) {
	query := chi.URLParam(r, "query")
	if query == "" {
		common.Fail(w, errors.New("Must include search query"), http.StatusBadRequest)
		return
	}
	if len([]byte(query)) > 256 {
		common.Fail(w, errors.New("Search query too long"), http.StatusBadRequest)
		return
	}
	results, err := models.RecommendedQuery(query)
	if err != nil {
		common.Fail(w, fmt.Errorf("Could not execute search: %s", err), http.StatusInternalServerError)
		return
	}
	common.JSON(w, map[string]interface{}{"recommended": results})
}
