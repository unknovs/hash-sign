package env

import "os"

var (
	PemFile = os.Getenv("PEM_FILE")
	ApiKey  = os.Getenv("API_KEY")
	AuthCert = os.Getenv("AUTH_CERT")
	SigningCert = os.Getenv("SIGN_CERT")
)