package main

import (
	"synthomat/minimi/internal"
	"synthomat/minimi/internal/db"
)

func main() {
	db := db.NewDB()

	internal.RunServer(db)
}
