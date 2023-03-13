package env

import "os"

var (
	PemFile = os.Getenv("PEM_FILE")
	ApiKey  = os.Getenv("API_KEY")
)