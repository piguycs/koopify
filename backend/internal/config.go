package internal

import (
	"crypto/tls"
	"errors"
	"strings"
)

const (
	DefaultPgDb     = "postgres://postgres:postgres@localhost:5432/?sslmode=disable"
	DefaultHostAddr = ":8080"
)

type Config struct {
	JwtSecret            string
	PgDb                 string
	AdyenApiKey          string
	AdyenMerchantAccount string
	AdyenThemeId         string
	CheckoutReturnUrl    string
	HostAddr             string
	TlsConfig            *tls.Config
}

func LoadConfig() (Config, error) {
	pgdb := GetEnvDefault("PGDB", DefaultPgDb)
	hostAddr := GetEnvDefault("HOST_ADDR", DefaultHostAddr)

	jwtSecret := GetEnv("JWT_SECRET")
	adyenApiKey := GetEnv("ADYEN_API_KEY")
	adyenMerchantAccount := GetEnv("ADYEN_MERCHANT_ACCOUNT")
	adyenThemeId := GetEnv("ADYEN_THEME_ID")
	checkoutReturnUrl := GetEnv("CHECKOUT_RETURN_URL")

	// go through all, make a list of the ones which are missing and notify the user
	missing := []string{}
	if jwtSecret == nil {
		missing = append(missing, "JWT_SECRET")
	}
	if adyenApiKey == nil {
		missing = append(missing, "ADYEN_API_KEY")
	}
	if adyenMerchantAccount == nil {
		missing = append(missing, "ADYEN_MERCHANT_ACCOUNT")
	}
	if adyenThemeId == nil {
		missing = append(missing, "ADYEN_THEME_ID")
	}
	if checkoutReturnUrl == nil {
		missing = append(missing, "CHECKOUT_RETURN_URL")
	}

	if len(missing) > 0 {
		return Config{}, errors.New("missing required environment variables: " + strings.Join(missing, ", "))
	}

	// pass in the environment variable keys for reading the TLS config
	tlsConfig, err := TlsConfig("TLS_ENABLED", "TLS_CERT", "TLS_KEY")
	if err != nil {
		return Config{}, err
	}

	return Config{
		JwtSecret:            *jwtSecret,
		PgDb:                 pgdb,
		AdyenApiKey:          *adyenApiKey,
		AdyenMerchantAccount: *adyenMerchantAccount,
		AdyenThemeId:         *adyenThemeId,
		CheckoutReturnUrl:    *checkoutReturnUrl,
		HostAddr:             hostAddr,
		TlsConfig:            tlsConfig,
	}, nil
}
