package main

import (
	"log"

	"github.com/ghofuobudgeghikhdrhiehdihebd/pc-usbc-int/internal/serialmanager"
)

func main() {
	manager := serialmanager.New()

	if err := manager.Connect(); err != nil {
		log.Printf("AirLock connection unavailable: %v", err)
	} else {
		defer manager.Disconnect()
		log.Println("Connected to AirLock")
	}

	runChatServer()
}
