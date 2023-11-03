package db

import (
	"time"

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