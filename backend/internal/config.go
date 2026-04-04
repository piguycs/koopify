package internal

import (
	"crypto/tls"
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

	var jwtSecret, adyenApiKey, adyenMerchantAccount, adyenThemeId, checkoutReturnUrl string

	err := GetEnvs([]EnvVarDef{
		{Value: &jwtSecret, Var: "JWT_SECRET"},
		{Value: &adyenApiKey, Var: "ADYEN_API_KEY"},
		{Value: &adyenMerchantAccount, Var: "ADYEN_MERCHANT_ACCOUNT"},
		{Value: &adyenThemeId, Var: "ADYEN_THEME_ID"},
		{Value: &checkoutReturnUrl, Var: "CHECKOUT_RETURN_URL"},
	})

	if err != nil {
		return Config{}, err
	}

	// pass in the environment variable keys for reading the TLS config
	tlsConfig, err := TlsConfig("TLS_ENABLED", "TLS_CERT", "TLS_KEY")
	if err != nil {
		return Config{}, err
	}

	return Config{
		JwtSecret:            jwtSecret,
		PgDb:                 pgdb,
		AdyenApiKey:          adyenApiKey,
		AdyenMerchantAccount: adyenMerchantAccount,
		AdyenThemeId:         adyenThemeId,
		CheckoutReturnUrl:    checkoutReturnUrl,
		HostAddr:             hostAddr,
		TlsConfig:            tlsConfig,
	}, nil
}
