package main

import (
	"encoding/base64"
	"time"

	"github.com/gofiber/fiber/v2"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Client struct {
	gorm.Model
	Name       string
	DocumentID string `gorm:"unique"`
	Licenses   []License
}

type License struct {
	gorm.Model
	UUID     string `gorm:"unique"`
	Expire   time.Time
	Delay    int
	ClientID uint
}

func openDB() *gorm.DB {
	dsn := "host=localhost user=gorm password=gorm dbname=gorm port=5432 sslmode=disable TimeZone=Europe/Madrid"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	db.AutoMigrate(&Client{}, &License{})
	if err != nil {
		panic(err)
	}
	return db
}

func createClient(db *gorm.DB) {
	oneLicense := License{
		UUID:   "98f4465b-d1ea-6786-5313-f9e13144f25d",
		Expire: time.Date(2024, 10, 21, 0, 0, 0, 0, time.FixedZone("Europe/Madrid", 0)),
		Delay:  0,
	}
	client := Client{
		Name:       "Hotel Torremolinos Centro",
		DocumentID: "H2988",
		Licenses: []License{
			oneLicense,
		},
	}
	db.Create(&client)
	db.Save(&client)
}

func getLicense(db *gorm.DB, key string) License {
	license := License{UUID: key}
	db.Select("id", "uuid", "expire", "delay").Find(&license)
	return license
}

func addLicense(db *gorm.DB, documentId string, licenseKey string, expireDate time.Time) License {
	client := Client{DocumentID: documentId}
	db.Find(&client)
	license := License{UUID: licenseKey, Delay: 0, Expire: expireDate}
	client.Licenses = append(client.Licenses, license)
	db.Save(&client)
	return license
}

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

	createClient(db)

	app.Get("/status/:uuid", func(c *fiber.Ctx) error {
		uuidb64 := c.Params("uuid")
		uuid := parseBase64(uuidb64)
		license := getLicense(db, uuid)
		return c.JSON(&fiber.Map{
			license.UUID: &fiber.Map{
				"id":     license.ID,
				"expire": license.Expire.Format("02.01.2006"),
				"delay":  license.Delay,
			},
		})

	})

	app.Get("/activate/:documentId/:uuid", func(c *fiber.Ctx) error {
		uuidb64 := c.Params("uuid")
		uuid := parseBase64(uuidb64)
		didb64 := c.Params("documentId")
		did := parseBase64(didb64)
		license := addLicense(db, did, uuid, time.Now().AddDate(1, 0, 0))
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
