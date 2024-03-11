package requests

type Request struct {
	EmptyAsice  string       `json:"emptyAsice"`
	SignedFiles []SignedFile `json:"signedFiles"`
}

type SignedFile struct {
	FileName    string `json:"fileName"`
	EncodedFile string `json:"encodedFile"`
}
