package model

import "time"

type Url struct {
	Id         int64     `db:"id"`
	Url        string    `db:"url"`
	Code       string    `db:"code"`
	CreatedAt  time.Time `db:"created_at"`
	ClickCount int64     `db:"click_count"`
}
