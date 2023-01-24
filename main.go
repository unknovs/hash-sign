package main

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
)

type hashSignature struct {
	Hash           string `json:"hash"`
	SignatureValue string `json:"signatureValue"`
}

var (
	certFile   = os.Getenv("PEM_FILE")
)

func main() {
	// Read the pem file
	certData, err := os.Open(certFile)
	if err != nil {
		log.Fatalf("Failed to read pem file: %s", err)
	}
	fmt.Println("certData readed")

	// Decode the pem file

	pemfileinfo, _ := certData.Stat()
	var size int64 = pemfileinfo.Size()
	pembytes := make([]byte, size)

	buffer := bufio.NewReader(certData)
	_, err = buffer.Read(pembytes)

	block, _ := pem.Decode([]byte(pembytes))
	if block == nil {
		log.Fatalf("Failed to decode pem file")
	}
	fmt.Println("PEM file decoded - size : ", size)
	certData.Close()

	// Extract the private key from the pem file
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		log.Fatalf("Failed to parse pem file: %s", err)
	}
	fmt.Println("Private Key loaded and ready to work!!!")

	// Router
	http.HandleFunc("/sign", func(w http.ResponseWriter, r *http.Request) {
		// Limit to POST only
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		var hashSignatureValue hashSignature
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

		// lets print some movement
		fmt.Println("Hash value: ", hashSignatureValue.Hash, " processed")

		// Write the JSON response
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(hashSignatureValue)
		if err != nil {
			log.Printf("Failed to encode JSON: %s", err)
			http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
		}
	},
	)

	log.Fatal(http.ListenAndServe(":80", nil))
}
