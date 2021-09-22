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
