package responses

type HashSignature struct {
	SignatureMethod string `json:"signatureMethod"`
	Hash            string `json:"hash"`
	SignatureValue  string `json:"signatureValue"`
}
