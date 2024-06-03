package functions

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"

	"github.com/unknovs/hash-sign/routes/requests"
	"github.com/unknovs/hash-sign/routes/responses"
)

func CalculateVerificationCode(w http.ResponseWriter, r *http.Request) {

	if !isPostMethod(r) {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	var req requests.RequestVerificationCode
	err = json.Unmarshal(body, &req)
	if err != nil {
		http.Error(w, "Error unmarshalling request body", http.StatusBadRequest)
		return
	}

	decodedHash, err := base64.StdEncoding.DecodeString(req.Hash)
	if err != nil {
		http.Error(w, "Error decoding base64 hash", http.StatusBadRequest)
		return
	}

	h := sha256.New()
	h.Write(decodedHash)
	hash := h.Sum(nil)

	lastTwoBytes := hash[len(hash)-2:]
	integer := int(lastTwoBytes[0])*256 + int(lastTwoBytes[1])
	verificationCode := integer % 10000

	res := responses.VerificationCodeResponse{
		VerificationCode: verificationCode,
	}

	jsonRes, err := json.Marshal(res)
	if err != nil {
		http.Error(w, "Error marshalling response body", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonRes)
}
