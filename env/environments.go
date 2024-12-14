package env

import (
	"log"
	"os"
	"strings"
)

var (
	PemFile          = os.Getenv("PEM_FILE")
	EcPemFile        = os.Getenv("EC_PEM_FILE")
	ApiKey           = os.Getenv("API_KEY")
	RsaAuthCert      = os.Getenv("RSA_AUTH_CERT")
	RsaSigningCert   = os.Getenv("RSA_SIGN_CERT")
	EcdsaAuthCert    = os.Getenv("ECDSA_AUTH_CERT")
	EcdsaSigningCert = os.Getenv("ECDSA_SIGN_CERT")
	jwtSigningKey    = os.Getenv("JWT_SIGNING_KEY")
)

// getEnvOrSecret reads the environment variable or Docker secret file content.
func getEnvOrSecret(varName string) string {
	value := os.Getenv(varName)
	if strings.HasPrefix(value, "/run/secrets/") {
		data, err := os.ReadFile(value)
		if err != nil {
			log.Printf("unable to read secret file %s: %v", value, err)
			return ""
		}
		return strings.TrimSpace(string(data))
	}
	return value
}
