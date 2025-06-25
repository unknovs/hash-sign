// SPDX-License-Identifier: MIT

// Copyright (c) 2024 Gatis Beikerts
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/unknovs/hash-sign/env"
	"github.com/unknovs/hash-sign/functions"
)

func main() {

	// Check if certificates exist and are valid
	err := functions.CheckCertificates()
	if err != nil {
		log.Printf("Failed to validate certificates: %s", err)
	}

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
	http.HandleFunc("/digest/verify", functions.APIKeyAuthorization(functions.VerifySignature))
	http.HandleFunc("/digest/calculateSummary", functions.APIKeyAuthorization(functions.HandleDigest))
	http.HandleFunc("/certificates", functions.APIKeyAuthorization(functions.HandleCertificatesRequest))
	http.HandleFunc("/asice/addFile", functions.APIKeyAuthorization(functions.HandleAddFileToAsiceRequest))
	http.HandleFunc("/encrypt/publicKey", functions.APIKeyAuthorization(functions.EncryptWithPublicKeyHandler))
	http.HandleFunc("/digest/verificationCode", functions.APIKeyAuthorization(functions.CalculateVerificationCode))
	http.HandleFunc("/jwt/generate", functions.APIKeyAuthorization(functions.JwtGenerateHandler))

	// Add a handler for the root path
	http.HandleFunc("/", functions.APIKeyAuthorization(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))

	fmt.Println("Server listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
