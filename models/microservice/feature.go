package microservice

import (
	"fmt"
	"time"

	"github.com/pigeatgarlic/ideacrawler/microservice/models/user"
)

type Feature struct {
	ID   int
	Name string
	Tags map[string]string

	Authority string

	EndpointIDs []int
	Allowed []user.Role
}

type CronJob struct {
	ID   int
	Name string
	Tags map[string]string

	Authority string

	Work []Endpoint

	Interval time.Duration
}

type Job struct {
	ID   int
	Name string
	Tags map[string]string

	Authority string

	Work []Endpoint
}

type Trigger struct {
	ID   int
	Name string
	Tags map[string]string

	Command string
	Feature Feature 
	Role user.Role
	Service MicroService
}


type ActionSerires map[int]*struct{
	EndpointID int
	MicroserviceID int
}

func (action ActionSerires) Add (endpoint Endpoint) error {
	action[endpoint.Order] = &struct {
		EndpointID int;
		MicroserviceID int;
	} {
		EndpointID: endpoint.ID,
		MicroserviceID : endpoint.MicroserviceID,
	}

	for i := 0; i < len(action); i++ {
		if action[i] == nil {
			return fmt.Errorf("invalid action, missing %dth action",i)
		}
	}
	return nil;
}