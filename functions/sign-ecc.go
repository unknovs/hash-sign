package functions

import (
	"crypto/ecdsa"
	"net/http"

	"github.com/unknovs/hash-sign/routes/requests"
)

func SigningHandlerEC(privateKey *ecdsa.PrivateKey) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !validateRequest(w, r, privateKey) {
			return
		}

		var signEcdsa requests.SignEcdsa
		if !decodeJSON(w, r, &signEcdsa) {
			return
		}

		hashBytes, err := decodeHash(signEcdsa)
		if err != nil {
			http.Error(w, "Failed to decode hash from base64", http.StatusBadRequest)
			return
		}

		signatureR, signatureS, err := signHash(privateKey, hashBytes)
		if err != nil {
			http.Error(w, "Error signing hash", http.StatusInternalServerError)
			return
		}

		signatureMethod := getSignatureMethod(r)
		signature, err := encodeSignature(signatureMethod, signatureR, signatureS, privateKey)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		sendResponse(w, signEcdsa, signature, signatureMethod)
	}
}
