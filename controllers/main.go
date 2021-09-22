package controllers

import (
	"github.com/gofiber/fiber/v2"
)

func Index(c *fiber.Ctx) error {
	data := fiber.Map{ "Title": "Hello, World!", }
	return c.Render("index", data)
}

