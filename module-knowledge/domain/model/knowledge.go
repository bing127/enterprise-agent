package model

type Knowledge struct {
    ID          string `json:"id"`
    Title       string `json:"title"`
    Description string `json:"description"`
    Content     string `json:"content"`
    CreatedAt   int64  `json:"created_at"`
    UpdatedAt   int64  `json:"updated_at"`
}