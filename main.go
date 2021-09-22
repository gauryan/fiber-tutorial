package main

import (
	"github.com/gauryan/xyz/routes"
	"github.com/gauryan/xyz/database"
)


func main() {
	app := routes.Router()
	database.Init()
	app.Listen(":3000")
}
