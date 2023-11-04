package db

import (
	"license_server/src/config"
	"license_server/src/utils"
	"log"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	conn *gorm.DB
}

func New(config config.Config) *DB {
	db, err := gorm.Open(postgres.Open(config.Dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}
	err = db.AutoMigrate(&Client{}, &License{})
	if err != nil {
		log.Fatalln(err)
	}
	return &DB{conn: db}
}

func (db *DB) CreateClient(name string, documentId string) {
	// oneLicense := License{
	// 	UUID:   "98f4465b-d1ea-6786-5313-f9e13144f25d",
	// 	Expire: time.Date(2024, 10, 21, 0, 0, 0, 0, time.FixedZone("Europe/Madrid", 0)),
	// 	Delay:  0,
	// }
	client := Client{
		Name:       name,
		DocumentID: documentId,
		Licenses: []License{},
	}
	db.conn.Create(&client)
	db.conn.Save(&client)
}

func (db *DB) getLicense(key string) (License, error) {
	license := License{UUID: key}
	tx := db.conn.Where(license).First(&license)
	return license, tx.Error
}

func (db *DB) AddLicense(documentId string, licenseKey string, expireDate time.Time) License {
	client := Client{DocumentID: documentId}
	db.conn.Find(&client)
	license := License{UUID: licenseKey, Delay: 0, Expire: expireDate}
	client.Licenses = append(client.Licenses, license)
	db.conn.Save(&client)
	return license
}

func (db *DB) GetLicenseStatus(c *fiber.Ctx) error {
	uuidb64 := c.Params("uuid")
	uuid, err := utils.ParseBase64(uuidb64)
	if (err != nil) {
		return c.SendStatus(http.StatusBadRequest)
	}
	license, err := db.getLicense(uuid)
	if (err != nil) {
		return c.SendStatus(http.StatusNotFound)
	}
	return c.JSON(&fiber.Map{
		license.UUID: &fiber.Map{
			"id":     license.ID,
			"expire": license.Expire.Format("02.01.2006"),
			"delay":  license.Delay,
		},
	})
}