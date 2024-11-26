package responses

type HashSignature struct {
	SessionId       string `json:"sessionId"`
	SignatureMethod string `json:"signatureMethod"`
	Hash            string `json:"hash"`
	SignatureValue  string `json:"signatureValue"`
}
