package microservice

import "fmt"

type MicroserviceListenerConfig struct {
	ClientID int

	MicroserviceID int
	Resource       string
	Namespace      string
}

func (config MicroserviceListenerConfig) GetTopic() string {
	return fmt.Sprintf("%s.%s.%d", config.Namespace, config.Resource, config.MicroserviceID)
}