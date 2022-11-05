package datastore_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/synthia-telemed/push-notification-consumer/pkg/datastore"
	pg "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"math/rand"
)

var _ = Describe("Patient DataStore", Ordered, func() {
	var (
		db                    *gorm.DB
		notificationDataStore datastore.NotificationDataStore
		patients              []*datastore.Patient
	)

	BeforeAll(func() {
		var err error
		db, err = gorm.Open(pg.Open(postgres.Config.DSN()), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
		Expect(err).To(BeNil())
	})

	BeforeEach(func() {
		rand.Seed(GinkgoRandomSeed())
		notificationDataStore = datastore.NewGormNotificationDataStore(db)

		Expect(db.AutoMigrate(&datastore.Patient{}, &datastore.Notification{})).To(Succeed())
		patients, _ = generatePatientWithNotifications(3)
		Expect(db.Create(&patients).Error).To(Succeed())
	})

	AfterEach(func() {
		Expect(db.Migrator().DropTable(&datastore.Patient{}, &datastore.Notification{})).To(Succeed())
	})

	Context("Create notification", func() {
		It("should crate notification", func() {
			noti := generateNotification(patients[1].ID)
			Expect(notificationDataStore.Create(&noti)).To(Succeed())
			var foundNoti datastore.Notification
			Expect(db.Where(&noti).First(&foundNoti).Error).To(BeNil())
		})
	})
})
