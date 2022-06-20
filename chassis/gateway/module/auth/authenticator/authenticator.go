package authenticator

import (
	"github.com/pigeatgarlic/goedf/chassis/util/config"
	usermodel "github.com/pigeatgarlic/goedf/models/user"
)

type Adapter struct {
	validateUrl string
}

func New(conf *config.SecurityConfig) *Adapter {
	var database Adapter
	database.validateUrl = conf.ValidatorUrl
	return &database
}

func (db *Adapter) ValidateUserRole(user *usermodel.User, role string) bool {
	return false
}
