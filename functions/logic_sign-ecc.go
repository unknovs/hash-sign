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

	"github.com/unknovs/hash-sign/routes/requests"
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

func decodeHash(signEcdsa requests.SignEcdsa) ([]byte, error) {
	// Decode the hash from base64
	hashBytes, err := base64.StdEncoding.DecodeString(signEcdsa.DigestToSign)
	if err != nil {
		log.Printf("Failed to decode hash from base64: %s", err)
		return nil, err
	}
	return hashBytes, nil
}

func signHash(privateKey *ecdsa.PrivateKey, hashBytes []byte) (*big.Int, *big.Int, error) {
	// Sign the hash
	signatureR, signatureS, err := ecdsa.Sign(rand.Reader, privateKey, hashBytes)
	if err != nil {
		log.Printf("Error signing hash: %s", err)
		return nil, nil, err
	}
	return signatureR, signatureS, nil
}

func encodeSignature(signatureMethod string, signatureR, signatureS *big.Int, privateKey *ecdsa.PrivateKey) ([]byte, error) {
	var signature []byte
	var err error

	switch signatureMethod {
	case "DER":
		// Marshal the signature to ASN.1 DER format
		signature, err = asn1.Marshal(struct {
			R, S *big.Int
		}{signatureR, signatureS})
		if err != nil {
			log.Printf("Failed to marshal ECDSA signature: %s", err)
			return nil, err
		}
	case "P1363":
		// Convert the R and S values to byte slices
		rBytes := signatureR.Bytes()
		sBytes := signatureS.Bytes()

		// Ensure the byte slices have the same length by prepending zeros if necessary
		keyBytes := (privateKey.Params().BitSize + 7) >> 3
		if len(rBytes) < keyBytes {
			temp := make([]byte, keyBytes)
			copy(temp[keyBytes-len(rBytes):], rBytes)
			rBytes = temp
		}
		if len(sBytes) < keyBytes {
			temp := make([]byte, keyBytes)
			copy(temp[keyBytes-len(sBytes):], sBytes)
			sBytes = temp
		}

		// Concatenate the R and S values
		signature = append(rBytes, sBytes...)
	default:
		return nil, fmt.Errorf("invalid signature method, use 'P1363' or 'DER'")
	}

	return signature, nil
}

func validateRequest(w http.ResponseWriter, r *http.Request, privateKey *ecdsa.PrivateKey) bool {
	if !isPostMethod(r) {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return false
	}

	if privateKey == nil {
		http.Error(w, "ECC Private key not loaded", http.StatusNotFound)
		return false
	}

	return true
}

func decodeJSON(w http.ResponseWriter, r *http.Request, signEcdsa *requests.SignEcdsa) bool {
	err := json.NewDecoder(r.Body).Decode(signEcdsa)
	if err != nil {
		log.Printf("Failed to decode JSON: %s", err)
		http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
		return false
	}

	return true
}

func getSignatureMethod(r *http.Request) string {
	signatureMethod := r.URL.Query().Get("SignatureMethod")
	if signatureMethod == "" {
		signatureMethod = r.URL.Query().Get("signatureMethod")
	}
	if signatureMethod == "" {
		signatureMethod = "DER"
	}

	return signatureMethod
}

func sendResponse(w http.ResponseWriter, signEcdsa requests.SignEcdsa, signature []byte, signatureMethod string) {
	signatureValue := base64.StdEncoding.EncodeToString(signature)

	hashSignature := responses.HashSignature{
		SignatureMethod: signatureMethod,
		Hash:            signEcdsa.DigestToSign,
		SignatureValue:  signatureValue,
	}

	log.Printf("Hash value: %v signed using %s method", hashSignature.Hash, signatureMethod)

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(hashSignature)
	if err != nil {
		log.Printf("Failed to encode JSON: %s", err)
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
	}
}
