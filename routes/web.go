package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	// "github.com/gauryan/xyz/controllers"
	"github.com/gauryan/xyz/controllers/mgmt"
)

func Router() *fiber.App {
	// App 생성과 템플릿 설정
	app := fiber.New(fiber.Config{
		Views: html.New("./views", ".html"),
	})

	// Route 설정
	// app.Get("/", controllers.Index)
	mgmtApp := app.Group("/mgmt")
	mgmtApp.Get("/admin", mgmt.ListAdmin)
	mgmtApp.Get("/admin/insert_form", mgmt.InsertForm)
	mgmtApp.Post("/admin/insert", mgmt.Insert)

	return app
}
