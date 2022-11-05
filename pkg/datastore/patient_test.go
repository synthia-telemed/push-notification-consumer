package datastore_test

import (
	"fmt"
	"github.com/google/uuid"
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
		db               *gorm.DB
		patientDataStore datastore.PatientDataStore
		patients         []*datastore.Patient
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
		var err error
		patientDataStore = datastore.NewGormPatientDataStore(db)
		Expect(err).To(BeNil())

		patients = generatePatients(5)
		Expect(db.AutoMigrate(&datastore.Patient{})).To(Succeed())
		Expect(db.Create(&patients).Error).To(BeNil())
	})

	AfterEach(func() {
		Expect(db.Migrator().DropTable(&datastore.Patient{})).To(Succeed())
	})

	Context("Find by ID or RefID", func() {
		When("found by ID", func() {
			It("should return patient with no error", func() {
				p := patients[2]
				patient, err := patientDataStore.FindByIDOrRefID(fmt.Sprintf("%d", p.ID))
				Expect(err).To(BeNil())
				Expect(patient.ID).To(Equal(p.ID))
				Expect(patient.RefID).To(Equal(p.RefID))
			})
		})

		When("found by RefID", func() {
			It("should return patient with no error", func() {
				p := patients[0]
				patient, err := patientDataStore.FindByIDOrRefID(p.RefID)
				Expect(err).To(BeNil())
				Expect(patient.ID).To(Equal(p.ID))
				Expect(patient.RefID).To(Equal(p.RefID))
			})
		})

		When("not found", func() {
			It("should return patient with no error", func() {
				patient, err := patientDataStore.FindByIDOrRefID(uuid.NewString())
				Expect(err).To(BeNil())
				Expect(patient).To(BeNil())
			})
		})
	})
})
