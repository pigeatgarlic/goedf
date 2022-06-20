package authenticator

import (
	usermodel "github.com/pigeatgarlic/goedf/models/user"
)

type Adapter interface {
	ValidateUserRole(user *usermodel.User, role string) bool
}
