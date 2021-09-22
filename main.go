package main

import (
	"github.com/gauryan/xyz/routes"
)


func main() {
	app := routes.Router()
	app.Listen(":3000")
}
