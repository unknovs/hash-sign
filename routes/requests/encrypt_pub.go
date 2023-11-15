package requests

type EncryptRequest struct {
	PublicKey     string `json:"public_key"`
	DataToEncrypt string `json:"dataToEncrypt"`
}