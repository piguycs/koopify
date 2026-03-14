package internal

import (
	"crypto/tls"
	"errors"
)

const (
	DefaultPgDb     = "postgres://postgres:postgres@localhost:5432/?sslmode=disable"
	DefaultHostAddr = ":8080"
	InvalidSecret   = "INVALID_SECRET"
)

type Config struct {
	JwtSecret string
	PgDb      string
	HostAddr  string
	TlsConfig *tls.Config
}

func LoadConfig() (Config, error) {
	jwtSecret := GetEnvDefault("JWT_SECRET", InvalidSecret)
	if jwtSecret == InvalidSecret {
		return Config{}, errors.New("cannot run the application without a valid JWT secret")
	}

	pgdb := GetEnvDefault("PGDB", DefaultPgDb)
	hostAddr := GetEnvDefault("HOST_ADDR", DefaultHostAddr)

	// pass in the environment variable keys for reading the TLS config
	tlsConfig, err := TlsConfig("TLS_ENABLED", "TLS_CERT", "TLS_KEY")
	if err != nil {
		return Config{}, err
	}

	return Config{
		JwtSecret: jwtSecret,
		PgDb:      pgdb,
		HostAddr:  hostAddr,
		TlsConfig: tlsConfig,
	}, nil
}
