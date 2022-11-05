package datastore_test

import (
	"github.com/google/uuid"
	"github.com/synthia-telemed/push-notification-consumer/pkg/datastore"
	"math/rand"
)

func generatePatient() *datastore.Patient {
	return &datastore.Patient{
		RefID: uuid.New().String(),
	}
}

func generatePatients(num int) []*datastore.Patient {
	users := make([]*datastore.Patient, num)
	for i := 0; i < num; i++ {
		users[i] = generatePatient()
	}
	return users
}

func generateNotification(patientID uint) datastore.Notification {
	return datastore.Notification{
		Title:     uuid.NewString(),
		Body:      uuid.NewString(),
		IsRead:    rand.Float32() > 0.5,
		PatientID: patientID,
	}
}

func generateNotifications(patientID uint, n int) ([]datastore.Notification, int) {
	notifications := make([]datastore.Notification, n, n)
	readCount := 0
	for i := 0; i < n; i++ {
		notifications[i] = generateNotification(patientID)
		if notifications[i].IsRead {
			readCount++
		}
	}
	return notifications, readCount
}

func generatePatientWithNotifications(n int) ([]*datastore.Patient, []int) {
	patients := make([]*datastore.Patient, n, n)
	readCounts := make([]int, n, n)
	for i := 0; i < n; i++ {
		patients[i] = generatePatient()
		patients[i].Notification, readCounts[i] = generateNotifications(patients[i].ID, rand.Intn(20)+2)
	}
	return patients, readCounts
}
