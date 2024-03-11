package functions

import (
	"crypto"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/unknovs/hash-sign/routes/requests"
)

func VerifySignature(w http.ResponseWriter, r *http.Request) {
	var verifyBody requests.VerifyBody
	err := json.NewDecoder(r.Body).Decode(&verifyBody)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse request body: %v", err), http.StatusUnprocessableEntity)
		return
	}

	certificateBytes, err := base64.StdEncoding.DecodeString(verifyBody.Certificate)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid certificate: %v", err), http.StatusBadRequest)
		return
	}

	certificate, err := x509.ParseCertificate(certificateBytes)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse certificate: %v", err), http.StatusBadRequest)
		return
	}

	publicKey, ok := certificate.PublicKey.(*rsa.PublicKey)
	if !ok {
		http.Error(w, fmt.Sprintf("Invalid public key: %T", certificate.PublicKey), http.StatusBadRequest)
		return
	}

	signatureBytes, err := base64.StdEncoding.DecodeString(verifyBody.SignatureValue)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid signature value: %v", err), http.StatusBadRequest)
		return
	}

	digestValue, err := base64.StdEncoding.DecodeString(verifyBody.DigestValue)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid digest value: %v", err), http.StatusBadRequest)
		return
	}

	err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, digestValue, signatureBytes)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to verify signature: %v", err), http.StatusBadRequest)
		return
	}

	fmt.Fprintln(w, "Signature is valid!")
}
