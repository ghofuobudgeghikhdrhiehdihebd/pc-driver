package mcpcomm

import (
	"encoding/json"
	"net/http"

	statepkg "github.com/ghofuobudgeghikhdrhiehdihebd/pc-usbc-int/internal/state"
	"github.com/ghofuobudgeghikhdrhiehdihebd/pc-usbc-int/pkg/httphelpers"
)

func RegisterHandler(state *statepkg.State) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			location string
		}
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			httphelpers.WriteError(w, 500, "failed to read json: "+err.Error())
			return
		}
		state.McpUdsLoc = request.location
		w.WriteHeader(200)
		return
	}
}
