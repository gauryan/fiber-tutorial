package mgmt

// controllers/mgmt

import (
	"github.com/gauryan/xyz/database"
	"github.com/gofiber/fiber/v2"
)

// Admin 목록
func ListAdmin(c *fiber.Ctx) error {
	type Admin struct {
		Sno    int
		Userid string
		Nick   string
	}
	var admins []Admin

	db := database.DBConn
	// db.Raw("SELECT sno, userid, nick FROM admins").Scan(&admins)
	db.Raw("CALL listAdmins()").Scan(&admins)

	data := fiber.Map{"Admins": admins}
	return c.Render("mgmt/admin/index", data, "mgmt/base")
}


// 관리자 추가 폼
func InsertForm(c *fiber.Ctx) error {
	return c.Render("mgmt/admin/insert_form", fiber.Map{})
}


// 관리자 추가
func Insert (c *fiber.Ctx) error {
	userid  := c.FormValue("userid")
	passwd1 := c.FormValue("passwd1")
	passwd2 := c.FormValue("passwd2")
	nick    := c.FormValue("nick")

	if passwd1 != passwd2 {
		return c.Redirect("/mgmt/admin")
	}
	db := database.DBConn
	db.Exec("CALL insertAdmin(?, ?, ?)", userid, passwd1, nick)

	return c.Redirect("/mgmt/admin")
}
