package main

import (
	"encoding/base64"
	"time"

	"github.com/gofiber/fiber/v2"
)

func parseBase64(uuid string) string {
	rawDecodedText, err := base64.StdEncoding.DecodeString(uuid)
	if err != nil {
		panic(err)
	}
	return string(rawDecodedText)
}

func main() {
	app := fiber.New()
	db := openDB()

	app.Get("/status/:uuid", db.getLicenseStatus)

	app.Get("/activate/:uuid", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	app.Post("/activate/:documentId/:uuid", func(c *fiber.Ctx) error {
		uuidb64 := c.Params("uuid")
		uuid := parseBase64(uuidb64)
		didb64 := c.Params("documentId")
		did := parseBase64(didb64)
		license := db.addLicense(did, uuid, time.Now().AddDate(1, 0, 0))
		return c.JSON(&fiber.Map{
			license.UUID: &fiber.Map{
				"id":     license.ID,
				"expire": license.Expire.Format("02.01.2006"),
				"delay":  license.Delay,
			},
		})
	})

	app.Listen(":3000")
}
