package main

import (
	"log"
	"synthomat/minimi/internal"
	"synthomat/minimi/internal/db"
)

func main() {
	baseConfig := internal.NewDefaultConfig()

	if baseConfig.AutoSecret {
		log.Printf("Using generated password: %s", baseConfig.AuthSecret)
	}

	gdb := db.NewDB(baseConfig.DBFileName)

	server := internal.NewServer(baseConfig)
	server.RunServer()
}
