package functions

import (
	"crypto"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/unknovs/hash-sign/env"
	"github.com/unknovs/hash-sign/routes/requests"
)

func VerifyHandler(privateKey *rsa.PrivateKey) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Limit to POST only
		if !isPostMethod(r) {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}
		// Parse request body
		verifyBody, err := parseVerifyBody(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// Verify signature
		err = verifySignature(verifyBody, privateKey)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		log.Printf("Signature in signatureValue is valid!")

		// If signature is validated, respond with this message
		message := "Signature is valid!"
		fmt.Fprintln(w, message)
	}
}

func parseVerifyBody(r *http.Request) (*requests.VerifyBody, error) {
	var verifyBody requests.VerifyBody
	err := json.NewDecoder(r.Body).Decode(&verifyBody)
	if err != nil {
		return nil, err
	}
	return &verifyBody, nil
}

func verifySignature(verifyBody *requests.VerifyBody, privateKey *rsa.PrivateKey) error {
	// Decode base64 certificate
	certificateBytes, err := base64.StdEncoding.DecodeString(verifyBody.Certificate)
	if err != nil {
		return errors.New("invalid certificate")
	}
	// Parse certificate
	certificate, err := parseCertificate(certificateBytes)
	if err != nil {
		return err
	}
	// Get public key from certificate
	publicKey, ok := certificate.PublicKey.(*rsa.PublicKey)
	if !ok {
		return errors.New("invalid public key")
	}
	// Decode received base64 signatureValue
	signatureBytes, err := base64.StdEncoding.DecodeString(verifyBody.SignatureValue)
	if err != nil {
		return errors.New("invalid signature value")
	}
	// Decode received base64 digestValue
	digestValue, err := base64.StdEncoding.DecodeString(verifyBody.DigestValue)
	if err != nil {
		return errors.New("invalid digest value")
	}
	// Verify signature using PKCS1v15
	err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, digestValue, signatureBytes)
	if err != nil {
		return errors.New("failed to verify signature")
	}
	return nil
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
	VerifyHandler(privateKey)(w, r)
}
