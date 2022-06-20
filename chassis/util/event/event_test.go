package event_testing

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/pigeatgarlic/goedf/chassis/util/event/provider/kafka"

	"github.com/pigeatgarlic/goedf/chassis/util/config"
	eslogger "github.com/pigeatgarlic/goedf/chassis/util/logger/es"
	"github.com/pigeatgarlic/goedf/models/event"
	"github.com/pigeatgarlic/goedf/models/microservice"
)

func TestKafkaMQ(t *testing.T) {
	os.Setenv("ENV", "unittest")
	defer func() {
		os.Setenv("ENV", "")
	}()

	conf := &config.ESLogConfig{
		ESurl:        "http://elasticsearch-1653540188.elasticstack.svc.cluster.local:9200",
		WarningIndex: "testwarning",
		InforIndex:   "testinfor",
		ErrorIndex:   "testerror",
		HostName:     "unittest",
		StdLog:       "true",
		Namespace:    "unittest",
	}

	var err error
	log, _ := eslogger.InitLogger(conf)

	mq := config.KafkaConfig{
		Server:       "cp-helm-charts-1653713557-cp-kafka.kafka.svc.cluster.local:9092",
		ClientID:     "0",
		Acks:         "0",
		DataProtocol: "protobuf",
	}

	querier, err := kafka.InitKafkaQuerier(log, &mq)
	if err != nil {
		t.Error(err)
		return
	}
	pusher, err := kafka.InitKafkaPusher(log, &mq)
	if err != nil {
		t.Error(err)
		return
	}

	querier.ConfigureTopic(microservice.MicroserviceListenerConfig{
		ClientID:       0,
		MicroserviceID: 0,
		Resource:       "unit",
		Namespace:      "test",
	})

	pusher.ConfigTopic([]microservice.MicroserviceListenerConfig{microservice.MicroserviceListenerConfig{
		ClientID:       0,
		MicroserviceID: 0,
		Resource:       "unit",
		Namespace:      "test",
	}})

	id := int(time.Now().UnixNano())
	log.Infor(fmt.Sprintf("Initialize test with event ID %d", id))
	eve := event.Event{
		ID:      id,
		Headers: map[string]string{"test": "true"},
		Actions: []event.Action{event.Action{
			ID: 0,

			Prev: 0,
			Next: 0,

			Service:  0,
			Endpoint: 0,

			Done:            true,
			SignedAuthority: make([]string, 0),

			Result: event.Result{
				Data:  map[string]string{},
				Error: "",
			},
		},
		},
	}

	err = pusher.Push(&eve)
	if err != nil {
		t.Error(err)
		return
	}

	doneChan := make(chan bool)
	go func() {
		for {
			result := querier.WaitIncomingEvent()
			if result == nil {
				t.Error("null event")
				return
			}
			if result.ID == eve.ID {
				log.Infor(fmt.Sprintf("Test case match event ID %d with result event ID %d", result.ID, eve.ID))
				doneChan <- true
			}
		}
	}()
	<-doneChan
}
