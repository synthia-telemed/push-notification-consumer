package config

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"github.com/synthia-telemed/push-notification-consumer/pkg/datastore"
	"github.com/synthia-telemed/push-notification-consumer/pkg/notification"
)

type Config struct {
	Notification notification.Config
	RabbitMQ     RabbitMQ
	SentryDSN    string `env:"SENTRY_DSN"`
	DB           datastore.Config
}

type RabbitMQ struct {
	User         string `env:"RABBITMQ_USER" envDefault:"guest"`
	Password     string `env:"RABBITMQ_PASSWORD" envDefault:"guest"`
	Host         string `env:"RABBITMQ_HOST" envDefault:"localhost"`
	Port         string `env:"RABBITMQ_PORT" envDefault:"5672"`
	QueueName    string `env:"RABBITMQ_QUEUE_NAME" envDefault:"push-notification-queue"`
	ExchangeName string `env:"RABBITMQ_EXCHANGE_NAME" envDefault:"notification"`
	RoutingKey   string `env:"RABBITMQ_ROUTING_KEY" envDefault:"push-notification"`
}

func (r RabbitMQ) GetURL() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%s", r.User, r.Password, r.Host, r.Port)
}

func Load() (*Config, error) {
	_ = godotenv.Load()
	cfg := &Config{}
	err := env.Parse(cfg)
	return cfg, err
}
