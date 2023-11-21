package domain

import "time"

type Url struct {
	Id        int       `db:"id" json:"id"`
	SourceUrl string    `db:"source_url" json:"source_url"`
	ShortUrl  string    `db:"short_url" json:"short_url"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
