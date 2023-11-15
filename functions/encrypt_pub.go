package functions

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"log"
	"net/http"

	"github.com/unknovs/hash-sign/routes/requests"
	"github.com/unknovs/hash-sign/routes/responses"
)

func EncryptWithPublicKeyHandler(w http.ResponseWriter, r *http.Request) {

	if !isPostMethod(r) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var inputData requests.EncryptRequest
	err := json.NewDecoder(r.Body).Decode(&inputData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	publicKey, err := GetPublicKey(inputData.PublicKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	encryptedData, err := EncryptWithPublicKey([]byte(inputData.DataToEncrypt), publicKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	encryptedDataResponse := responses.EncryptResponse{EncryptedData: base64.StdEncoding.EncodeToString(encryptedData)}
	jsonResponse, err := json.Marshal(encryptedDataResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsonResponse)
	if err != nil {
		log.Printf("Failed to write a response: %s", err)
		http.Error(w, "Failed to write a response", http.StatusInternalServerError)
	}
}

func GetPublicKey(pemStr string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte("-----BEGIN PUBLIC KEY-----\n" + pemStr + "\n-----END PUBLIC KEY-----"))
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	publicKey, ok := publicKeyInterface.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("failed to cast public key to RSA public key")
	}

	return publicKey, nil
}

func EncryptWithPublicKey(data []byte, publicKey *rsa.PublicKey) ([]byte, error) {
	encryptedData, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, data)
	if err != nil {
		return nil, err
	}

	return encryptedData, nil
}
