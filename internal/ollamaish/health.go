package ollamaish

import (
	"net/http"

	"github.com/ghofuobudgeghikhdrhiehdihebd/pc-usbc-int/pkg/httphelpers"
)

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httphelpers.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	httphelpers.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
