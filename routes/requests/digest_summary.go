package requests

type DigestSummaryRequest struct {
	DigestToSign string `json:"digest"`
}
