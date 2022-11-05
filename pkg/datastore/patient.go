package datastore

import (
	"errors"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type Patient struct {
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index"`
	RefID             string         `json:"refID" gorm:"unique"`
	ID                uint           `json:"id" gorm:"autoIncrement,primaryKey"`
	Notification      []Notification `gorm:"foreignKey:PatientID"`
	NotificationToken string         `json:"-"`
}

type PatientDataStore interface {
	FindByIDOrRefID(cred string) (*Patient, error)
}

type GormPatientDataStore struct {
	db *gorm.DB
}

func NewGormPatientDataStore(db *gorm.DB) *GormPatientDataStore {
	return &GormPatientDataStore{db: db}
}

func (d GormPatientDataStore) FindByIDOrRefID(id string) (*Patient, error) {
	var patient Patient
	q := d.db.Model(&Patient{})
	uintID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		q = q.Where("ref_id = ?", id)
	} else {
		q = q.Where("id = ?", uintID)
	}

	if err := q.First(&patient).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &patient, nil
}
