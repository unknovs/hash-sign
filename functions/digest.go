package functions

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func HandleDigest(w http.ResponseWriter, r *http.Request) {
	if !isGetMethod(r) {
		w.WriteHeader(http.StatusMethodNotAllowed)
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

	// Hash binary digest with SHA256
	hashedDigest := sha256.Sum256(binaryDigest)

	// Hex to base64
	digestSummary := base64.StdEncoding.EncodeToString(hashedDigest[:])

	// Convert digestSummary to JSON
	jsonData := map[string]string{"digestSummary": digestSummary}
	jsonBytes, err := json.Marshal(jsonData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error marshaling JSON: %v", err)
		return
	}

	log.Printf("Digest summary of digest %v calculated", digest)

	// Set headers and write response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	var writeErr error
	_, err = w.Write(jsonBytes)
	if err != nil {
		log.Printf("Digest summary of digest %v calculated", writeErr)
	}

}
