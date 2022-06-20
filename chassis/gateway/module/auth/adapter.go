package authenticator

import (
	usermodel "github.com/pigeatgarlic/ideacrawler/microservice/models/user"
)

type Adapter interface {
	ValidateUserRole(user *usermodel.User, role string) bool 
}