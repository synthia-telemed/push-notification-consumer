package datastore

import "fmt"

type Config struct {
	Host        string `env:"DATABASE_HOST" envDefault:"localhost"`
	TimeZone    string `env:"DATABASE_TIMEZONE" envDefault:"Asia/Bangkok"`
	User        string `env:"DATABASE_USER" envDefault:"postgres"`
	Password    string `env:"DATABASE_PASSWORD" envDefault:"postgres"`
	Name        string `env:"DATABASE_NAME,required" envDefault:"synthia"`
	SSLMode     string `env:"DATABASE_SSL_MODE" envDefault:"disable"`
	SSLRootCert string `env:"DATABASE_SSL_ROOT_CERT" envDefault:""`
	Port        int    `env:"DATABASE_PORT" envDefault:"5432"`
}

func (c Config) DSN() string {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%v sslmode=%s TimeZone=%s",
		c.Host, c.User, c.Password, c.Name, c.Port, c.SSLMode, c.TimeZone)
	if len(c.SSLRootCert) != 0 {
		dsn = fmt.Sprintf("%s sslrootcert=%s", dsn, c.SSLRootCert)
	}
	return dsn
}
