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

func SigningHandler(privateKey *rsa.PrivateKey) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Limit to POST only

		if !isPostMethod(r) {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		// Check if the private key is loaded
		if privateKey == nil {
			http.Error(w, "RSA Private key not loaded", http.StatusNotFound)
			return
		}

		var hashSignatureRequests []struct {
			SessionId string `json:"sessionId"`
			Hash      string `json:"hash"`
		}

		// Decode the incoming JSON array
		err := json.NewDecoder(r.Body).Decode(&hashSignatureRequests)
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
