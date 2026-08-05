package model

import "time"

type URL struct {
	ID         int64     `db:"id"`
	URL        string    `db:"url"`
	Code       string    `db:"code"`
	CreatedAt  time.Time `db:"created_at"`
	ClickCount int64     `db:"click_count"`
}
