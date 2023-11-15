package responses

type CertificatesResponse struct {
	AuthenticationCertificate string `json:"authentication_certificate,omitempty"`
	SigningCertificate        string `json:"signing_certificate,omitempty"`
}