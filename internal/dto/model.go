package dto

import "time"

type AddNewUrlRequests struct {
	Url string `json:"url"`
}

type GetCodeRequests struct {
	Url string `json:"url"`
}

type CodeResponse struct {
	Code string `json:"url"`
}

type UrlStatsResponse struct {
	Url        string    `json:"url"`
	Code       string    `json:"code"`
	CreatedAt  time.Time `json:"created_at"`
	ClickCount int64     `json:"click_count"`
}
