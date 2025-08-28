package models

import "time"

// DTOs
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

type PostUpdateDTO struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

// Filters
type PostFilters struct {
	ShortContent bool `json:"shortContent"`
	Limit        int  `json:"limit"`
	Page         int  `json:"page"`
}

// Render
type PostPageData struct {
	Posts    []Post
	NextPage int
}
