package main

import (
	"encoding/base64"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	app.Get("/activate/:uuid", func(c *fiber.Ctx) error {
		uuidb64 := c.Params("uuid")
		rawDecodedText, err := base64.StdEncoding.DecodeString(uuidb64)
		if err != nil {
			panic(err)
		}
		uuid := string(rawDecodedText)

		return c.JSON(&fiber.Map{
			uuid: &fiber.Map{
				"id":     1,
				"expire": "21.10.2024",
				"delay":  0,
			},
		})
	})

	app.Listen(":3000")
}
