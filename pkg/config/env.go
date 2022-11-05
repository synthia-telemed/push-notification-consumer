package config

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"github.com/synthia-telemed/push-notification-consumer/pkg/notification"
)

type Config struct {
	Notification notification.Config
	RabbitMQ     RabbmitMQ
}

type RabbmitMQ struct {
	User     string `env:"RABBITMQ_USER" envDefault:"guest"`
	Password string `env:"RABBITMQ_PASSWORD" envDefault:"guest"`
	Host     string `env:"RABBITMQ_HOST" envDefault:"localhost"`
	Port     string `env:"RABBITMQ_PORT" envDefault:"5672"`
}

func (r RabbmitMQ) GetURL() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%s", r.User, r.Password, r.Host, r.Port)
}

func Load() (*Config, error) {
	_ = godotenv.Load()
	cfg := &Config{}
	err := env.Parse(cfg)
	return cfg, err
}