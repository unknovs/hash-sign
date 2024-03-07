package functions

import (
	"bufio"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/x509"
	"encoding/asn1"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"

	"github.com/unknovs/hash-sign/routes/responses"
)

func GetECPrivateKey(filename string) (*ecdsa.PrivateKey, error) {
	//Lets check, if PEM file is added to env
	if filename == "" {
		log.Println("ECC Private key is not loaded.")
		return nil, nil
	}

	// Read the PEM file
	keyFile, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open ECC private key file: %w", err)
	}
	defer keyFile.Close()

	// Read the PEM file content
	keyFileInfo, _ := keyFile.Stat()
	keySize := keyFileInfo.Size()
	keyBytes := make([]byte, keySize)

	buffer := bufio.NewReader(keyFile)
	_, err = buffer.Read(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to read ECC private key file: %w", err)
	}

	// Decode PEM block
	block, _ := pem.Decode(keyBytes)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	// Parse PKCS#8 private key
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ECC PKCS#8 private key: %w", err)
	}

	// Convert to ECDSA private key
	ecKey, ok := key.(*ecdsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("parsed key is not ECDSA private key")
	}

	return ecKey, nil
}

func SigningHandlerEC(privateKey *ecdsa.PrivateKey) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Limit to POST only
		if !isPostMethod(r) {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		// Check if the private key is loaded
		if privateKey == nil {
			http.Error(w, "ECC Private key not loaded", http.StatusNotFound)
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

		// Decode the hash from base64
		hashBytes, err := base64.StdEncoding.DecodeString(hashSignatureValue.Hash)
		if err != nil {
			log.Printf("Failed to decode hash from base64: %s", err)
			http.Error(w, "Failed to decode hash from base64", http.StatusBadRequest)
			return
		}

		// Sign the hash
		signatureR, signatureS, err := ecdsa.Sign(rand.Reader, privateKey, hashBytes)
		if err != nil {
			log.Printf("Error signing hash: %s", err)
			http.Error(w, "Error signing hash", http.StatusInternalServerError)
			return
		}

		// Marshal the signature to ASN.1 DER format
		ecdsaSignature, err := asn1.Marshal(struct {
			R, S *big.Int
		}{signatureR, signatureS})
		if err != nil {
			log.Printf("Failed to marshal ECDSA signature: %s", err)
			http.Error(w, "Failed to marshal ECDSA signature", http.StatusInternalServerError)
			return
		}

		// Encode the signature to base64
		signatureValue := base64.StdEncoding.EncodeToString(ecdsaSignature)
		hashSignatureValue.SignatureValue = signatureValue

		// Log the signed hash value
		log.Printf("Hash value: %v signed", hashSignatureValue.Hash)

		// Write the JSON response
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(hashSignatureValue)
		if err != nil {
			log.Printf("Failed to encode JSON: %s", err)
			http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
		}
	}
}
