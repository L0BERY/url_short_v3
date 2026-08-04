package dto

import "url_shortener_v3/internal/model"

func NewCodeResponse(url *model.Url) CodeResponse {
	return CodeResponse{
		Code: url.Code,
	}
}

func NewUrlStatsResponse(url *model.Url) UrlStatsResponse {
	return UrlStatsResponse{
		Url:        url.Url,
		Code:       url.Code,
		CreatedAt:  url.CreatedAt,
		ClickCount: url.ClickCount,
	}
}
