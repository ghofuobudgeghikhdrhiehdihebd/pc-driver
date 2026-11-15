package main

import (
	"net/http"
	"time"

	"github.com/ghofuobudgeghikhdrhiehdihebd/pc-usbc-int/internal/ollamaish"
	"github.com/ghofuobudgeghikhdrhiehdihebd/pc-usbc-int/pkg/httphelpers"
)

func main() {
	mux := http.NewServeMux()
	ollamaish.RegisterAll(mux)
	server := &http.Server{
		Addr:              ":11434",
		Handler:           httphelpers.WithCORS(httphelpers.WithLogging(mux)),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
