package main

import (
	"synthomat/minimi/internal"
	"synthomat/minimi/internal/db"
)

func main() {
	gdb := db.NewDB()

	internal.RunServer(gdb)
}
