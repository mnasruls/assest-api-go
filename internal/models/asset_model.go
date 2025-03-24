package models

import (
	"time"

	"gorm.io/gorm"
)

type Asset struct {
	Id              string     `json:"id"`
	Name            string     `json:"name"`
	Type            string     `json:"type"`
	Value           float64    `json:"value"`
	AcquisitionDate time.Time  `json:"acquisition_date"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at"`
}

func (a Asset) TableName() string {
	return "assets"
}

func (l *Asset) BeforeCreate(tx *gorm.DB) (err error) {
	tNow := time.Now().UTC()
	l.CreatedAt = tNow
	l.UpdatedAt = tNow
	return
}

func (l *Asset) BeforeUpdate(tx *gorm.DB) (err error) {
	l.UpdatedAt = time.Now().UTC()
	return
}
