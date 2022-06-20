package tenant

import (
	"github.com/pigeatgarlic/goedf/models/request-response/response"
	"github.com/pigeatgarlic/goedf/models/user"
)

type Tenant struct {
	UserID    int
	SessionID uint64

	Name  string
	Roles []user.Role
	Tags  map[string]string

	channel chan (*response.UserResponse)
}

func NewTenant(ID uint64, user *user.User) *Tenant {
	return &Tenant{
		UserID:    user.ID,
		SessionID: ID,

		Name:  user.UserName,
		Roles: user.Roles,

		channel: make(chan *response.UserResponse),
	}
}

func (tenant *Tenant) SendResponse(resp *response.UserResponse) {
	tenant.channel <- resp
}
func (tenant *Tenant) ListenonResponse() (resp *response.UserResponse) {
	return <-tenant.channel
}
func (tenant *Tenant) Terminate() {
	close(tenant.channel)
}
