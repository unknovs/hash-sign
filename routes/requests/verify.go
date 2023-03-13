package requests

type VerifyBody struct {
	SignatureValue string `json:"signatureValue"`
	Certificate    string `json:"certificate"`
	DigestValue    string `json:"digestValue"`
}