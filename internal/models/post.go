package models

import "time"

type Post struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type PostCreateDTO struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}
