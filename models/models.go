package models

type ShortenRequest struct {
	LongUrl string `json:"long_url"`
}

type ShortenResponse struct {
	ShortUrl string `json:"short_url"`
}
