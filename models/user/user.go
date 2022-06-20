package user

import "time"

type User struct {
	ID int

	UserName            string
	NormalLizedUserName string

	FullName            string
	NormalLizedFullName string

	Email            string
	NormalLizedEmail string
	EmailConfirmed   bool

	PasswordHash  string
	SecurityStamp string

	PhoneNumber          string
	PhoneNumberConfirmed bool

	Address string
	CreatedAt   time.Time
	DateOfBirth time.Time


	TwoFactorEnabled bool
	LoginFailedCount int

	IsValidated bool
	Roles []Role 					
}


type Role struct {
	ID          int
	Name        string
	Description string
	CreatedAt   time.Time
	Tags        map[string]string
	Users		[]User 				
}

type UserLogin struct {
	Failed bool
	Provider string
	IpAddress string
	CreatedAt   time.Time
}

