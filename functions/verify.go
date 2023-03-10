package functions

import (
	"crypto"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/unknovs/hash-sign/env"
	"github.com/unknovs/hash-sign/routes/requests"
)
func verifyHandler(privateKey *rsa.PrivateKey) http.HandlerFunc {
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

		// Parse request body
		var verifyBody requests.VerifyBody
		err := json.NewDecoder(r.Body).Decode(&verifyBody)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Decode base64 certificate
		certificateBytes, err := base64.StdEncoding.DecodeString(verifyBody.Certificate)
		if err != nil {
			http.Error(w, "Invalid certificate", http.StatusBadRequest)
			return
		}

		// Parse certificate
		certificate, err := parseCertificate(certificateBytes)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Get public key from certificate
		publicKey, ok := certificate.PublicKey.(*rsa.PublicKey)
		if !ok {
			http.Error(w, "Invalid public key", http.StatusBadRequest)
			return
		}

		// Decode received base64 signatureValue
		signatureBytes, err := base64.StdEncoding.DecodeString(verifyBody.SignatureValue)
		if err != nil {
			http.Error(w, "Invalid signature value", http.StatusBadRequest)
			return
		}

		// Decode received base64 digestValue
		digestValue, err := base64.StdEncoding.DecodeString(verifyBody.DigestValue)
		if err != nil {
			http.Error(w, "Invalid digest value", http.StatusBadRequest)
			return
		}

		// Verify signature using PKCS1v15
		err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, digestValue, signatureBytes)
		if err != nil {
			http.Error(w, "Failed to verify signature", http.StatusBadRequest)
			return
		}

		// If signature is validated, respond with this message
		message := "Signature is valid!"
		fmt.Fprintln(w, message)
	}
}

func parseCertificate(certificateBytes []byte) (*x509.Certificate, error) {
	certificate, err := x509.ParseCertificate(certificateBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse certificate: %v", err)
	}

	// Check that the public key algorithm is RSA
	if _, ok := certificate.PublicKey.(*rsa.PublicKey); !ok {
		return nil, fmt.Errorf("invalid public key algorithm: %T", certificate.PublicKey)
	}

	return certificate, nil
}

func VerifyHandlerWrapper(w http.ResponseWriter, r *http.Request) {
	privateKey, err := GetPrivateKey(env.PemFile)
	if err != nil {
		http.Error(w, "Failed to get private key", http.StatusBadRequest)
	}
	verifyHandler(privateKey)(w, r)
}

