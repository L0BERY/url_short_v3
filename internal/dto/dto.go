package dto

import "time"

type AddNewURLRequests struct {
	URL string `json:"url"`
}

type GetCodeRequests struct {
	URL string `json:"url"`
}

type CodeResponse struct {
	Code string `json:"code"`
}

type URLStatsResponse struct {
	URL        string    `json:"url"`
	Code       string    `json:"code"`
	CreatedAt  time.Time `json:"created_at"`
	ClickCount int64     `json:"click_count"`
}
