package functions

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const logMessage = "%s?key=%s&type=%s responded"

func writeError(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	fmt.Fprintf(w, `{"error": "%s"}`, message)
}

func HandleCertificatesRequest(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	key := query.Get("key")
	certType := query.Get("type")

	if !isGetMethod(r) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if !certificatesExist() {
		writeError(w, http.StatusNotFound, "No certificates found in environment")
		return
	}

	response, err := getCertificatesResponse(key, certType, r)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	jsonBytes, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error marshaling JSON: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonBytes)
	if err != nil {
		log.Printf("Error writing JSON response: %v", err)
	}
}
