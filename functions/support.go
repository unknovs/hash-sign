package functions

import (
	"log"
	"net/http"
	"os"

	"github.com/unknovs/hash-sign/env"
)

func isPostMethod(r *http.Request) bool {
	return r.Method == http.MethodPost
}

func isGetMethod(r *http.Request) bool {
	return r.Method == http.MethodGet
}

func APIKeyAuthorization(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if env.ApiKey != "" && env.ApiKey != r.Header.Get("API-Key") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}

func CheckVolumeMounted(volumePath string) bool {
	// Attempt to obtain information about the volume without creating a directory
	if _, err := os.Stat(volumePath); err != nil {
		if os.IsNotExist(err) {
			log.Printf("Volume is not mounted at %s", volumePath)
		} else {
			log.Printf("Error checking volume: %v", err)
		}
		return false
	}

	log.Printf("Volume is mounted at %s", volumePath)
	return true
}
