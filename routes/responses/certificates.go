package responses

type CertificatesResponse struct {
	RsaAuthenticationCertificate   string `json:"rsa_authentication_certificate,omitempty"`
	RsaSigningCertificate          string `json:"rsa_signing_certificate,omitempty"`
	EcdsaAuthenticationCertificate string `json:"ecdsa_authentication_certificate,omitempty"`
	EcdsaSigningCertificate        string `json:"ecdsa_signing_certificate,omitempty"`
}
