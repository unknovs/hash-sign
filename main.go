package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/unknovs/hash-sign/env"
	"github.com/unknovs/hash-sign/functions"
)

func main() {
    
	// Check if API key is set and shall be used
	if env.ApiKey == "" {
		log.Println("API key not set in environment. Continuing without API key")
	}

	// Check if the volume is mounted
	volumePath := "/tmp" 
	if !functions.CheckVolumeMounted(volumePath) {
		log.Println("Volume is not available or mounted. asice/addFile method wont be available")
	}

	// Read the PEM file and extract the private key
	privateKey, err := functions.GetPrivateKey(env.PemFile)
	if err != nil {
		log.Fatalf("Failed to parse pem file: %s", err)
	}

	// Router
	http.HandleFunc("/digest/sign", functions.APIKeyAuthorization(functions.SigningHandler(privateKey)))
	http.HandleFunc("/digest/verify", functions.APIKeyAuthorization(functions.VerifyHandlerWrapper))
	http.HandleFunc("/digest/calculateSummary/", functions.APIKeyAuthorization(functions.HandleDigest))
	http.HandleFunc("/certificates", functions.APIKeyAuthorization(functions.HandleCertificatesRequest))
	http.HandleFunc("/asice/addFile", functions.APIKeyAuthorization(functions.HandleAddFileToAsiceRequest))
	http.HandleFunc("/encrypt/publicKey", functions.APIKeyAuthorization(functions.EncryptWithPublicKeyHandler))

	fmt.Println("Server listening on port 80...")
	log.Fatal(http.ListenAndServe(":80", nil))
}
