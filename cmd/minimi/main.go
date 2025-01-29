package main

import (
	"synthomat/minimi/internal"
)

func main() {
	db := internal.NewDB()

	internal.RunServer(db)

}
