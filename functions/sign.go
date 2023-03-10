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
		log.Println("Private key is not loaded.")
		return nil, nil
	}

	// Read the pem file
	certData, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read pem file: %w", err)
	}
	defer certData.Close()

	pemfileinfo, _ := certData.Stat()
	var size int64 = pemfileinfo.Size()
	pembytes := make([]byte, size)

	buffer := bufio.NewReader(certData)
	_, err = buffer.Read(pembytes)
	if err != nil {
		return nil, fmt.Errorf("failed to read pem file: %w", err)
	}

	block, _ := pem.Decode([]byte(pembytes))
	if block == nil {
		return nil, fmt.Errorf("failed to decode pem file")
	}

	// Extract the private key from the pem file
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse pem file: %w", err)
	}

	fmt.Println("Private Key loaded and ready to work!!!")

	return privateKey, nil
}

func SigningHandler(privateKey *rsa.PrivateKey) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Limit to POST only
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		// Check API key
		if !checkAPIKey(r) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Check if the private key is loaded
		if privateKey == nil {
			http.Error(w, "Private key not loaded", http.StatusNotFound)
			return
		}

		var hashSignatureValue responses.HashSignature
		
		// Decode the incoming JSON
		err := json.NewDecoder(r.Body).Decode(&hashSignatureValue)
		if err != nil {
			log.Printf("Failed to decode JSON: %s", err)
			http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
			return
		}

		// Sign the hash
		hashBytes, _ := base64.StdEncoding.DecodeString(hashSignatureValue.Hash)
		signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashBytes[:])
		if err != nil {
			log.Printf("Error signing hash: %s", err)
			http.Error(w, "Error signing hash", http.StatusInternalServerError)
			return
		}

		// Encode the signature to base64
		signatureValue := base64.StdEncoding.EncodeToString(signature)
		hashSignatureValue.SignatureValue = signatureValue
		hashSignatureValue.Hash = base64.StdEncoding.EncodeToString(hashBytes)

		// Log the signed hash value
		fmt.Println("Hash value: ", hashSignatureValue.Hash, " processed")

		// Write the JSON response
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(hashSignatureValue)
		if err != nil {
			log.Printf("Failed to encode JSON: %s", err)
			http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
		}
	}
}
