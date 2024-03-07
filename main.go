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
		log.Printf("Failed to parse RSA pem file: %s", err)
		log.Printf("signing using RSA key /digest/sign wont be possible")
	}

	// Read the ECC PEM file and extract the private key
	privateKeyEC, err := functions.GetECPrivateKey(env.EcPemFile)
	if err != nil {
		log.Printf("Failed to parse ECC pem file: %s", err)
		log.Printf("signing using ecc key /digest/sign-ecc wont be possible")
	}

	// Router
	http.HandleFunc("/digest/sign", functions.APIKeyAuthorization(functions.SigningHandler(privateKey)))
	http.HandleFunc("/digest/sign-ecc", functions.APIKeyAuthorization(functions.SigningHandlerEC(privateKeyEC)))
	http.HandleFunc("/digest/verify", functions.APIKeyAuthorization(functions.VerifyHandlerWrapper))
	http.HandleFunc("/digest/calculateSummary/", functions.APIKeyAuthorization(functions.HandleDigest))
	http.HandleFunc("/certificates", functions.APIKeyAuthorization(functions.HandleCertificatesRequest))
	http.HandleFunc("/asice/addFile", functions.APIKeyAuthorization(functions.HandleAddFileToAsiceRequest))
	http.HandleFunc("/encrypt/publicKey", functions.APIKeyAuthorization(functions.EncryptWithPublicKeyHandler))

	fmt.Println("Server listening on port 80...")
	log.Fatal(http.ListenAndServe(":80", nil))
}
