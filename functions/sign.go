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

package functions

import (
	"bufio"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/unknovs/hash-sign/routes/responses"
)

func GetPrivateKey(filename string) (*rsa.PrivateKey, error) {

	//Lets check, if PEM file is added to env
	if filename == "" {
		log.Println("RSA Private key is not loaded.")
		return nil, nil
	}

	// Read the pem file
	certData, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read RSA pem file: %w", err)
	}
	defer certData.Close()

	pemfileinfo, _ := certData.Stat()
	var size int64 = pemfileinfo.Size()
	pembytes := make([]byte, size)

	buffer := bufio.NewReader(certData)
	_, err = buffer.Read(pembytes)
	if err != nil {
		return nil, fmt.Errorf("failed to read RSA pem file: %w", err)
	}

	block, _ := pem.Decode([]byte(pembytes))
	if block == nil {
		return nil, fmt.Errorf("failed to decode RSA pem file")
	}

	// Extract the private key from the pem file
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse RSA pem file: %w", err)
	}

	return privateKey, nil
}

// Single request without sessionId
type SingleHashRequest struct {
	SessionId string `json:"sessionId,omitempty"`
	Hash      string `json:"hash"`
}

// Array request with sessionId
type HashSignatureRequest struct {
	SessionId string `json:"sessionId"`
	Hash      string `json:"hash"`
}

func SigningHandler(privateKey *rsa.PrivateKey) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !isPostMethod(r) {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		if privateKey == nil {
			http.Error(w, "RSA Private key not loaded", http.StatusNotFound)
			return
		}

		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}

		// Try to parse single request first
		var singleRequest SingleHashRequest
		var hashSignatureRequests []HashSignatureRequest
		err = json.Unmarshal(bodyBytes, &singleRequest)
		if err == nil && singleRequest.Hash != "" {
			// Single request handling
			hashBytes, _ := base64.StdEncoding.DecodeString(singleRequest.Hash)
			signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashBytes[:])
			if err != nil {
				log.Printf("Error signing hash: %s", err)
				http.Error(w, "Error signing hash", http.StatusInternalServerError)
				return
			}

			// Single response
			hashSignatureResponse := responses.HashSignature{
				SessionId:       singleRequest.SessionId,
				SignatureMethod: "PKCS1v15",
				Hash:            singleRequest.Hash,
				SignatureValue:  base64.StdEncoding.EncodeToString(signature),
			}

			log.Printf("Hash value: %v signed", hashSignatureResponse.Hash)

			w.Header().Set("Content-Type", "application/json")
			err = json.NewEncoder(w).Encode(hashSignatureResponse) // Note: no array here
			if err != nil {
				log.Printf("Failed to encode JSON: %s", err)
				http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
			}
		} else {
			// Try array format
			err = json.Unmarshal(bodyBytes, &hashSignatureRequests)
			if err != nil {
				log.Printf("Failed to decode JSON: %s", err)
				http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
				return
			}

			var hashSignatureResponses []responses.HashSignature

			// Process each hash in the array
			for _, request := range hashSignatureRequests {
				hashBytes, _ := base64.StdEncoding.DecodeString(request.Hash)
				signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashBytes[:])
				if err != nil {
					log.Printf("Error signing hash: %s", err)
					http.Error(w, "Error signing hash", http.StatusInternalServerError)
					return
				}

				hashSignatureResponse := responses.HashSignature{
					SessionId:       request.SessionId,
					SignatureMethod: "PKCS1v15",
					Hash:            request.Hash,
					SignatureValue:  base64.StdEncoding.EncodeToString(signature),
				}

				hashSignatureResponses = append(hashSignatureResponses, hashSignatureResponse)
			}

			// Log the signed hash values
			for _, response := range hashSignatureResponses {
				log.Printf("Hash value: %v signed", response.Hash)
			}

			// Write the JSON response
			w.Header().Set("Content-Type", "application/json")
			err = json.NewEncoder(w).Encode(hashSignatureResponses)
			if err != nil {
				log.Printf("Failed to encode JSON: %s", err)
				http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
			}
		}
	}
}
