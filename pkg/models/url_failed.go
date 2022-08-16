package models

type FailedURL struct {
	URL    string `json:"url"`
	Reason string `json:"reason"`
}
