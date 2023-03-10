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

	// Read the PEM file and extract the private key
	privateKey, err := functions.GetPrivateKey(env.PemFile)
	if err != nil {
		log.Fatalf("Failed to parse pem file: %s", err)
	}

	// Router
	http.HandleFunc("/sign", functions.SigningHandler(privateKey))
	http.HandleFunc("/verify", functions.VerifyHandlerWrapper)

	fmt.Println("Server listening on port 80...")
	log.Fatal(http.ListenAndServe(":80", nil))
}
