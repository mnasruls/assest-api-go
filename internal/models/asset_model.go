package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Asset struct {
	Id              string     `json:"id" gorm:"primaryKey;type:varchar(36);not null"`
	Name            string     `json:"name" gorm:"type:varchar(255);not null"`
	Type            string     `json:"type" gorm:"type:varchar(255);not null"`
	Value           float64    `json:"value" gorm:"type:float;not null"`
	AcquisitionDate time.Time  `json:"acquisition_date" gorm:"type:date;not null"`
	CreatedAt       time.Time  `json:"created_at" gorm:"type:timestamp;not null"`
	UpdatedAt       time.Time  `json:"updated_at" gorm:"type:timestamp;not null"`
	DeletedAt       *time.Time `json:"deleted_at" gorm:"type:timestamp;default:null"`
}

func (a Asset) TableName() string {
	return "assets"
}

func (l *Asset) BeforeCreate(tx *gorm.DB) (err error) {
	tNow := time.Now().UTC()
	if l.Id == "" {
		l.Id = uuid.New().String()
	}
	l.CreatedAt = tNow
	l.UpdatedAt = tNow
	return
}

func (l *Asset) BeforeUpdate(tx *gorm.DB) (err error) {
	l.UpdatedAt = time.Now().UTC()
	return
}
