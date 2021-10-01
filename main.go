package main

import (
	"github.com/gauryan/xyz/routes"
	"github.com/gauryan/xyz/database"
	"github.com/gauryan/xyz/store"
)


func main() {
	app := routes.Router()
	database.Init()
	store.Init()
	app.Listen(":3000")
}
