package phosphor

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/samuelhorwitz/phosphorescence/api/common"
	"github.com/samuelhorwitz/phosphorescence/api/models"
)

func Search(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	if query == "" {
		common.Fail(w, errors.New("Must include search query"), http.StatusBadRequest)
		return
	}
	if len([]byte(query)) > 256 {
		common.Fail(w, errors.New("Search query too long"), http.StatusBadRequest)
		return
	}
	result, err := models.Query(query)
	if err != nil {
		common.Fail(w, fmt.Errorf("Could not execute search: %s", err), http.StatusInternalServerError)
		return
	}
	common.JSON(w, map[string]interface{}{"result": result})
}
