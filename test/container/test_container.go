package container

import (
	"context"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
)

type Terminate func() error

type TestContainer struct {
	Host      string
	Port      *nat.Port
	Terminate Terminate
}

//type TestContainer interface {
//	Terminate(ctx context.Context) error
//}

func NewTestContainer(req testcontainers.ContainerRequest, port string) (*TestContainer, error) {
	ctx := context.Background()
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}
	host, err := container.Host(ctx)
	if err != nil {
		return nil, err
	}
	var p *nat.Port
	if port != "" {
		mappedPort, err := container.MappedPort(ctx, nat.Port(port))
		if err != nil {
			return nil, err
		}
		p = &mappedPort
	}

	return &TestContainer{
		Host: host,
		Port: p,
		Terminate: func() error {
			return container.Terminate(ctx)
		},
	}, nil
}
