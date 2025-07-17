package dto

import (
	"time"
)

type RegisterRequest struct {
	Login    string `json:"login" binding:"required,min=3,max=20,alphanum"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginRequest struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserResponse struct {
	ID    uint   `json:"id"`
	Login string `json:"login"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type CreateAdRequest struct {
	Title       string  `json:"title" binding:"required,min=5,max=100"`
	Description string  `json:"description" binding:"required,min=10,max=1000"`
	ImageURL    string  `json:"image_url" binding:"required,url"`
	Price       float64 `json:"price" binding:"required,gt=0"`
}

type AdvertisementResponse struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	ImageURL    string    `json:"image_url"`
	Price       float64   `json:"price"`
	CreatedAt   time.Time `json:"created_at"`
	AuthorLogin string    `json:"author_login"`
	IsOwner     bool      `json:"is_owner,omitempty"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type AdFilterRequest struct {
	Page      int     `form:"page" binding:"min=1"`
	Limit     int     `form:"limit" binding:"min=1,max=100"`
	SortBy    string  `form:"sort_by" binding:"oneof=created_at price"`
	SortOrder string  `form:"sort_order" binding:"oneof=asc desc"`
	MinPrice  float64 `form:"min_price" binding:"min=0"`
	MaxPrice  float64 `form:"max_price" binding:"min=0"`
}
