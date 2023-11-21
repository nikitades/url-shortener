package domain

import "time"

type Visit struct {
	Id        int       `db:"id" json:"id"`
	UserAgent string    `db:"user_agent" json:"user_agent"`
	UrlId     int       `db:"url_id" json:"url_id"`
	UrlSource string    `db:"url_source" json:"url_source"`
	UrlCode   string    `db:"url_code" json:"url_code"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
