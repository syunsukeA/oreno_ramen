package object

import (
	"time"
)

type User struct {
	UserID int64 `json:"user_id" db:"user_id"`
	UserName string `json:"username" db:"username"`
	EMail string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
	ProfileImg string `json:"profile_img" db:"profile_img"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}