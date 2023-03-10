package functions

import (
	"net/http"

	"github.com/unknovs/hash-sign/env"
)

func checkAPIKey(r *http.Request) bool {
	if env.ApiKey == "" {
		return true
	}

	receivedAPIKey := r.Header.Get("API-Key")
	return receivedAPIKey == env.ApiKey
}