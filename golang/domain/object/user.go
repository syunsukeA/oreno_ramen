package object

import (
	"time"
)

type User struct {
	UserID int64
	UserName string
	EMail string
	Password string
	ProfileImg string
	CreatedAt time.Time
	UpdatedAt time.Time
}