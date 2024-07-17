package responses

type DigestSummary struct {
	DigestSummary  string `json:"digestSummary"`
	UrLSafeSummary string `json:"URLSafeDigestSummary"`
	Algorithm      string `json:"algorithmUsed"`
}
