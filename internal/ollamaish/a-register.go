package ollamaish

import (
	statepkg "github.com/ghofuobudgeghikhdrhiehdihebd/pc-usbc-int/internal/state"
)

func RegisterAll(state *statepkg.State) {
	state.Mux.HandleFunc("/api/chat", ChatHandler)
	state.Mux.HandleFunc("/health", HealthHandler)
}
