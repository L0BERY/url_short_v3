package dto

import "url_shortener_v3/internal/model"

func NewCodeResponse(url *model.URL) CodeResponse {
	return CodeResponse{
		Code: url.Code,
	}
}

func NewUrlStatsResponse(url *model.URL) URLStatsResponse {
	return URLStatsResponse{
		URL:        url.URL,
		Code:       url.Code,
		CreatedAt:  url.CreatedAt,
		ClickCount: url.ClickCount,
	}
}
