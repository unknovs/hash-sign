package responses

type HashSignature struct {
	Hash           string `json:"hash"`
	SignatureValue string `json:"signatureValue"`
}