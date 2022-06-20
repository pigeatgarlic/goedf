package kafka

import (
	"encoding/json"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/pigeatgarlic/ideacrawler/microservice/chassis/util/config"
	"github.com/pigeatgarlic/ideacrawler/microservice/chassis/util/logger"
	"github.com/pigeatgarlic/ideacrawler/microservice/models/event"
	"github.com/pigeatgarlic/ideacrawler/microservice/models/microservice"
)

type KafkaListener struct {
	consumer *kafka.Consumer
	config   config.MQInterface

	run    bool

	logger logger.Logger

	eventQueue    chan (*event.Event)
	kafkaQueue    chan (kafka.Event)
}

func (querier *KafkaListener) handleKafkaMessage(message *kafka.Message) {
	var KafkaEvent *event.Event
	var err error

	switch querier.config.GetProtocol() {
	case "protobuf":
		KafkaEvent, err = event.FromProtobytes(message.Value)
	case "json":
		var eve event.Event
		err = json.Unmarshal(message.Value, &eve);
		KafkaEvent = &eve;
	}

	if err != nil {
		querier.logger.Error(fmt.Sprintf("Fail to decode kafka message: %s",err.Error()))
		return
	}

	querier.eventQueue <- KafkaEvent
}

func InitKafkaQuerier(logger logger.Logger,
	config config.MQInterface) (*KafkaListener, error) {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": 			   config.GetServer(),
		"group.id":          			   config.GetClientID(),
		"session.timeout.ms":              6000,
		"go.events.channel.enable":        true,
		"go.application.rebalance.enable": true,
		"auto.offset.reset":    		   "earliest",
	})
	if err != nil {
		return nil, err
	}

	ret := &KafkaListener{
		config:		   config,
		logger: 	   logger,	
		consumer: 	   consumer,
		eventQueue:    make(chan *event.Event),
		kafkaQueue:    nil,
		run:		   true,
	}

	return ret, nil
}

func (querier *KafkaListener) WaitIncomingEvent() *event.Event {
	if !querier.run  {
		return nil;
	}
	if querier.eventQueue == nil {
		querier.logger.Error("waiting for incoming event while unconfigured");
		return nil
	}
	event := <-querier.eventQueue;
	querier.logger.Infor(fmt.Sprintf("Handle incoming event %d",event.ID));
	return event
}

func (listener *KafkaListener) ConfigureTopic(topic microservice.MicroserviceListenerConfig) {
	listener.consumer.SubscribeTopics([]string{topic.GetTopic()}, nil)
	listener.kafkaQueue = listener.consumer.Events()

	listener.logger.Infor("Configuring event listener topics");
	go func() {
		for listener.run {
			ev := <-listener.kafkaQueue
			switch e := ev.(type) {
			case *kafka.Message:
				go listener.handleKafkaMessage(e)
			case kafka.Error:
				listener.logger.Error(fmt.Sprintf("Kafka error %v", e))
			case kafka.AssignedPartitions:
				listener.logger.Debug(fmt.Sprintf("Assigned partitions %v", e))
				listener.consumer.Assign(e.Partitions)
			case kafka.RevokedPartitions:
				listener.logger.Debug(fmt.Sprintf("Revoke partition %v", e))
				listener.consumer.Unassign()
			case kafka.PartitionEOF:
				fmt.Printf("%% Reached %v\n", e)
			default:
				listener.logger.Warning(fmt.Sprintf("Ignored %v", e))
			}
		}
	}()
}