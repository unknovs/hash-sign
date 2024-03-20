package functions

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/unknovs/hash-sign/routes/responses"
)

func HandleDigest(w http.ResponseWriter, r *http.Request) {
	if !isGetMethod(r) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Parse the query parameters
	query := r.URL.Query()
	hash := query.Get("hash")
	if hash == "" {
		hash = "sha256" // Default to sha256 if hash is not provided
	} else if hash != "sha256" && hash != "sha384" && hash != "sha512" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Unsupported hash algorithm: %s", hash)
		return
	}

	digest := r.URL.Path[len("/digest/calculateSummary/"):]

	// Base64 to binary
	binaryDigest, err := base64.StdEncoding.DecodeString(digest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error decoding base64 digest: %v", err)
		return
	}

	// Hash binary digest with selected algorithm
	var hashedDigest []byte
	var algorithm string
	switch hash {
	case "sha256":
		hashed := sha256.Sum256(binaryDigest)
		hashedDigest = hashed[:]
		algorithm = "sha256"
	case "sha384":
		hashed := sha512.Sum384(binaryDigest)
		hashedDigest = hashed[:]
		algorithm = "sha384"
	case "sha512":
		hashed := sha512.Sum512(binaryDigest)
		hashedDigest = hashed[:]
		algorithm = "sha512"
	}

	// Hex to base64
	digestSummary := base64.StdEncoding.EncodeToString(hashedDigest)

	response := responses.DigestSummary{
		DigestSummary: digestSummary,
		Algorithm:     algorithm,
	}

	jsonBytes, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error marshaling JSON: %v", err)
		return
	}

	// spam a bit in log
	log.Printf("Digest summary of digest %v calculated with hash %s", digest, hash)

	// Set headers and write response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	var writeErr error
	_, err = w.Write(jsonBytes)
	if err != nil {
		log.Printf("Error writing JSON response: %v", writeErr)
	}
}
