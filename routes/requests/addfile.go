package requests

type Request struct {
	EmptyAsice  string `json:"emptyAsice"`
	SignedFiles []File `json:"signedFiles"`
}

type File struct {
	FileName    string `json:"fileName"`
	EncodedFile string `json:"encodedFile"`
}
