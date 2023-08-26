package model

type Author struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	TwitterId   string `json:"twitter_id"`
	Email       string `json:"email" validate:"required"`
	PhoneNumber string `json:"phone_number" validate:"required"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}
