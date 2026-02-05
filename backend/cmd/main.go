package main

import (
	"context"
	"crypto/tls"
	"net/http"
	"os"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

var CommitHash = "dev"

func main() {
	e := echo.New()

	e.Use(middleware.RequestLogger())
	// this API should be available for standalone use
	e.Use(middleware.CORS("*"))

	e.GET("/commit", func(c *echo.Context) error {
		return c.String(http.StatusOK, CommitHash)
	})

	tlsConfig, err := TlsConfig()
	if err != nil {
		e.Logger.Error("failed to load TLS", "error", err)
	}

	sc := echo.StartConfig{Address: ":8080", TLSConfig: tlsConfig}
	if tlsConfig != nil {
		e.Logger.Info("starting server with TLS :)")
		sc.TLSConfig = tlsConfig
	} else {
		e.Logger.Info("starting server without TLS :(")
	}

	if err := sc.Start(context.Background(), e); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}

func TlsConfig() (*tls.Config, error) {
	if os.Getenv("TLS_ENABLED") == "" || os.Getenv("TLS_ENABLED") == "0" {
		return nil, nil
	}

	certPEM, err := os.ReadFile(os.Getenv("TLS_CERT"))
	if err != nil {
		return nil, err
	}
	keyPEM, err := os.ReadFile(os.Getenv("TLS_KEY"))
	if err != nil {
		return nil, err
	}

	cert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		return nil, err
	}

	return &tls.Config{Certificates: []tls.Certificate{cert}}, nil
}
