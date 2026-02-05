package main

import (
	"context"
	"net/http"
	"os"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

var CommitHash = "dev"

func main() {
	e := echo.New()

	e.Use(middleware.RequestLogger())

	e.GET("/commit", func(c *echo.Context) error {
		return c.String(http.StatusOK, CommitHash)
	})

	cert, key, err := TlsConfig()
	if err != nil {
		e.Logger.Error("failed to get tls config, probably TLS_CERT or TLS_KEY is unset", "error", err)
	}

	context := context.Background()
	sc := echo.StartConfig{
		Address: ":8080",
	}

	if cert != nil && key != nil {
		e.Logger.Info("starting server with TLS :)")
		if err := sc.StartTLS(context, e, *cert, *key); err != nil {
			e.Logger.Error("failed to start server", "error", err)
		}
	} else {
		e.Logger.Info("starting server without TLS :(")
		if err := sc.Start(context, e); err != nil {
			e.Logger.Error("failed to start server", "error", err)
		}
	}

}

func TlsConfig() (cert *string, key *string, err error) {
	tlsEnabled := os.Getenv("TLS_ENABLED")
	tlsCert := os.Getenv("TLS_CERT")
	tlsKey := os.Getenv("TLS_KEY")

	// TLS_ENABLED=1 means it is enabled, unset or 0 means it is not. Anything above 0 is true
	if tlsEnabled == "" || tlsEnabled == "0" {
		return nil, nil, nil
	}

	if tlsCert == "" || tlsKey == "" {
		return nil, nil, os.ErrNotExist
	}

	return &tlsCert, &tlsKey, nil
}
