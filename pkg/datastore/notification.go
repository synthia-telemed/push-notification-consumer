package datastore

import (
	"gorm.io/gorm"
	"time"
)

type Notification struct {
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	ID        uint           `json:"id" gorm:"autoIncrement,primaryKey"`
	Title     string         `json:"title"`
	Body      string         `json:"body"`
	IsRead    bool           `json:"is_read"`
	PatientID uint           `json:"patient_id"`
}

type NotificationDataStore interface {
	Create(notification *Notification) error
}

type GormNotificationDataStore struct {
	db *gorm.DB
}

func NewGormNotificationDataStore(db *gorm.DB) *GormNotificationDataStore {
	return &GormNotificationDataStore{db: db}
}

func (d GormNotificationDataStore) Create(notification *Notification) error {
	return d.db.Create(notification).Error
}
