package session

import (
	"time"

	"github.com/pigeatgarlic/goedf/models/user"
)

type Session struct {
	ID int

	User  user.User
	Login user.UserLogin

	CreatedAt time.Time
	ExpireAt  time.Time

	History map[time.Time]string
}
