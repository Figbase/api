package controllers

import "github.com/gofiber/fiber/v2"

func Home(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"error": false,
		"msg":   nil,
		"data": fiber.Map{
			"path": "/",
			"msg":  "Welcome, you've found the home endpoint",
		},
	})
}
