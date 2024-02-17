package main

import (
	"mygram-final/database"
	"mygram-final/routers"
)

func main() {
	database.StartDB()

	r := routers.StartApp()
	r.Run(":8080")
}
