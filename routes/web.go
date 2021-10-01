package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	"github.com/gauryan/xyz/store"
	"github.com/gauryan/xyz/controllers/mgmt"
)

// authMgmt 미들웨어
func authMgmt(c *fiber.Ctx) error {
	session, err := store.SessionStore.Get(c)
	if err != nil {
		panic(err)
	}

	mgmt_login := session.Get("mgmt-login")
	if (mgmt_login != true) {
		return c.Redirect("/mgmt")
	}

	return c.Next()
}

func Router() *fiber.App {
	// App 생성과 템플릿 설정
	app := fiber.New(fiber.Config{
		Views: html.New("./views", ".html"),
	})

	// Route 설정
	// app.Get("/", controllers.Index)
	mgmtApp1 := app.Group("/mgmt")
	mgmtApp1.Get("/", mgmt.Index)
	mgmtApp1.Post("/login", mgmt.Login)

	mgmtApp2 := app.Group("/mgmt", authMgmt)
	mgmtApp2.Get("/logout", mgmt.Logout)
	mgmtApp2.Get("/admin", mgmt.ListAdmin)
	mgmtApp2.Get("/admin/insert_form", mgmt.InsertForm)
	mgmtApp2.Post("/admin/insert", mgmt.Insert)
	mgmtApp2.Get("/admin/chg_passwd_form/:id", mgmt.ChgPasswdForm)
	mgmtApp2.Post("/admin/chg_passwd", mgmt.ChgPasswd)
	mgmtApp2.Get("/admin/update_form/:id", mgmt.UpdateForm)
	mgmtApp2.Post("/admin/update", mgmt.Update)
	mgmtApp2.Get("/admin/delete/:id", mgmt.Delete)

	return app
}
