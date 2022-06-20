package authenticator

import (
	"fmt"

	"github.com/pigeatgarlic/ideacrawler/microservice/chassis/gateway/module/auth/authenticator"
	"github.com/pigeatgarlic/ideacrawler/microservice/chassis/util/config"
	"github.com/pigeatgarlic/ideacrawler/microservice/models/user"
)

type Authenticator struct {
	config 	*config.SecurityConfig
	adapter Adapter
}

type jwtClaim struct {
	ID       int    `json:"ID"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func InitAuthenticator(config *config.SecurityConfig) *Authenticator {
	return &Authenticator{
		config: 	 config,
		adapter:     authenticator.New(config),
	}
}

func (auth *Authenticator) ValidateToken(signedToken string, role string) (*user.User, error) {
	valid := false
	var user user.User

	for i := 0; i < len(user.Roles); i++ {
		if user.Roles[i].Name == role {
			valid = true
		}
	}

	if valid {
		return &user, nil
	} else {
		return &user, fmt.Errorf("user do not have enough privillege")
	}
}
