package env

import "os"

var (
	PemFile          = os.Getenv("PEM_FILE")
	EcPemFile        = os.Getenv("EC_PEM_FILE")
	ApiKey           = os.Getenv("API_KEY")
	RsaAuthCert      = os.Getenv("RSA_AUTH_CERT")
	RsaSigningCert   = os.Getenv("RSA_SIGN_CERT")
	EcdsaAuthCert    = os.Getenv("ECDSA_AUTH_CERT")
	EcdsaSigningCert = os.Getenv("ECDSA_SIGN_CERT")
)
