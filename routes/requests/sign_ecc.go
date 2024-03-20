package requests

type SignEcdsa struct {
	DigestToSign string `json:"hash"`
}
