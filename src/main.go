package main

import (
	"license_server/src/db"
	"license_server/src/utils"
	"license_server/src/views"
	"time"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	db := db.New()

	app.Get("/status/:uuid", db.GetLicenseStatus)

	app.Get("/activate/:uuid", func(c *fiber.Ctx) error {
		uuid := utils.ParseBase64(c.Params("uuid"))
		index := views.Index(uuid)
		return utils.RenderResponse(c, index)
	})

	app.Post("/activate/:documentId/:uuid", func(c *fiber.Ctx) error {
		uuidb64 := c.Params("uuid")
		uuid := utils.ParseBase64(uuidb64)
		didb64 := c.Params("documentId")
		did := utils.ParseBase64(didb64)
		license := db.AddLicense(did, uuid, time.Now().AddDate(1, 0, 0))
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
