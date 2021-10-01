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
	mgmtApp.Get("/admin/chg_passwd_form/:id", mgmt.ChgPasswdForm)
	mgmtApp.Post("/admin/chg_passwd", mgmt.ChgPasswd)
	mgmtApp.Get("/admin/update_form/:id", mgmt.UpdateForm)
	mgmtApp.Post("/admin/update", mgmt.Update)
	mgmtApp.Get("/admin/delete/:id", mgmt.Delete)

	return app
}
