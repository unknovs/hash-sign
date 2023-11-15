package functions

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/unknovs/hash-sign/env"
	"github.com/unknovs/hash-sign/routes/responses"
)

func HandleCertificatesRequest(w http.ResponseWriter, r *http.Request) {
	var response responses.CertificatesResponse
	var err error

	if !isGetMethod(r) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if env.AuthCert == "" && env.SigningCert == "" {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "No certificates found in environment")
		return
	}

	if certType := r.URL.Query().Get("type"); certType == "auth" {
		response.AuthenticationCertificate = env.AuthCert
		log.Printf("%s?type=%s responded", r.URL.Path, certType)
	} else if certType == "sign" {
		response.SigningCertificate = env.SigningCert
		log.Printf("%s?type=%s responded", r.URL.Path, certType)
	} else if certType != "" {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Type used not allowed")
		return
	} else {
		response.AuthenticationCertificate = env.AuthCert
		response.SigningCertificate = env.SigningCert
		log.Printf("%s responded", r.URL.Path)
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error encoding response: %v", err)
	}
}
