package internal

import (
	"context"
	"crypto/tls"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
)

func TlsConfig(enabled_env, cert_env, privkey_env string) (*tls.Config, error) {
	if os.Getenv(enabled_env) == "" || os.Getenv(enabled_env) == "0" {
		return nil, nil
	}

	certPEM, err := os.ReadFile(os.Getenv(cert_env))
	if err != nil {
		return nil, err
	}
	keyPEM, err := os.ReadFile(os.Getenv(privkey_env))
	if err != nil {
		return nil, err
	}

	cert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		return nil, err
	}

	return &tls.Config{Certificates: []tls.Certificate{cert}}, nil
}

func GetEnvDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func StartServer(ctx context.Context, address string, tlsConfig *tls.Config, echoHandler *echo.Echo) {
	sc := echo.StartConfig{Address: address}
	if tlsConfig != nil {
		echoHandler.Logger.Info("starting server with TLS :)")
		sc.TLSConfig = tlsConfig
	} else {
		echoHandler.Logger.Info("starting server without TLS :(")
	}

	if err := sc.Start(ctx, echoHandler); err != nil {
		echoHandler.Logger.Error("failed to start server", "error", err)
	}
}

type CustomValidator struct {
	validator *validator.Validate
}

func NewCustomValidator() CustomValidator {
	return CustomValidator{validator: validator.New()}
}

func (cv *CustomValidator) Validate(i any) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.ErrBadRequest.Wrap(err)
	}
	return nil
}

// Bind a value from Echo's query parameters, path parameters or json body, and
// validate the provided values using reflection. Only return the underlying
// data if all the data exists AND it is validated
func BindAndValidate[T any](ctx *echo.Context) (*T, error) {
	data := new(T)

	if err := ctx.Bind(data); err != nil {
		return nil, err
	}

	if err := ctx.Validate(data); err != nil {
		return nil, err
	}

	return data, nil
}
