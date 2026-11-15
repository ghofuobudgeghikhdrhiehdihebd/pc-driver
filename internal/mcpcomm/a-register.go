package mcpcomm

import statepkg "github.com/ghofuobudgeghikhdrhiehdihebd/pc-usbc-int/internal/state"

func RegisterAll(state *statepkg.State) {
	state.Mux.HandleFunc("/mcpcomm/register", RegisterHandler(state))
}
