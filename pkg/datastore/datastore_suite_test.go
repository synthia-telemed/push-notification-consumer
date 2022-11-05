package datastore_test

import (
	"github.com/synthia-telemed/push-notification-consumer/pkg/datastore"
	"github.com/synthia-telemed/push-notification-consumer/test/container"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestDatastore(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Datastore Suite")
}

type PostgresDB struct {
	container.Terminate
	datastore.Config
}

var (
	postgres PostgresDB
)

var _ = BeforeSuite(func() {
	postgres = setupPostgresDBContainer()
})

var _ = AfterSuite(func() {
	Expect(postgres.Terminate()).To(Succeed())
})

func setupPostgresDBContainer() PostgresDB {
	config := datastore.Config{
		User:     "postgres",
		Password: "postgres",
		Name:     "synthia",
		SSLMode:  "disable",
		TimeZone: "Asia/Bangkok",
	}
	req := testcontainers.ContainerRequest{
		Image:        "postgres:13-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     config.User,
			"POSTGRES_PASSWORD": config.Password,
			"POSTGRES_DB":       config.Name,
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections").WithOccurrence(2),
	}
	con, err := container.NewTestContainer(req, "5432")
	Expect(err).To(BeNil())

	config.Host = con.Host
	config.Port = con.Port.Int()

	return PostgresDB{
		Config:    config,
		Terminate: con.Terminate,
	}
}
