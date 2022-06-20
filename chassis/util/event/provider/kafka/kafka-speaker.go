package kafka

// TODO handle failure during runtime
import (
	"encoding/json"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/pigeatgarlic/ideacrawler/microservice/chassis/util/config"
	"github.com/pigeatgarlic/ideacrawler/microservice/chassis/util/logger"
	"github.com/pigeatgarlic/ideacrawler/microservice/models/event"
	"github.com/pigeatgarlic/ideacrawler/microservice/models/microservice"
)

type KafkaSpeaker struct {
	resultChan chan (kafka.Event)
	channel    chan (*kafka.Message)

	producer   *kafka.Producer
	config     config.MQInterface

	topics		map[int]string
	run 		bool

	logger    	logger.Logger
}

func InitKafkaPusher(log logger.Logger,
	config config.MQInterface) (*KafkaSpeaker, error) {
	var err error

	producer,err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": config.GetServer(),
	})
	if err != nil {
		return nil, err
	}

	ret := KafkaSpeaker{
		run :		true,

		config:     config,
		logger:     log,

		producer: 	producer,
		topics:  	make(map[int]string),

		channel:	producer.ProduceChannel(),
		resultChan: producer.Events(),
	}

	go func() {
		for e := range ret.resultChan {
			switch ev := e.(type) {
			case *kafka.Message:
				m := ev
				if m.TopicPartition.Error != nil {
					ret.logger.Error(fmt.Sprintf("Delivery failed: %v\n", m.TopicPartition.Error))
				} else {
					ret.logger.Debug(fmt.Sprintf("Delivered message to topic %s [%d] at offset %v",
						*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset))
				}
				return

			default:
				ret.logger.Debug(fmt.Sprintf("Ignored event: %s\n", ev))
			}
		}
	}()



	return &ret, nil
}

func FreeKafkaPusher(pusher KafkaSpeaker) {
	close(pusher.channel);
	close(pusher.resultChan);
	pusher.producer.Close();
}

func (pusher *KafkaSpeaker) Push(event *event.Event) error {
	var data []byte
	var err error

	switch pusher.config.GetProtocol() {
	case "protobuf":
		data, err = event.ToProtobytes()
	case "json":
		data, err = json.Marshal(event)
	}

	if err != nil {
		pusher.logger.Error(fmt.Sprintf("Fail to decode kafka message %v", data))
		return err;
	}

	svcID := event.CurrentAction();
	if svcID == nil {
		err = fmt.Errorf("Invalid event %d: no action",event.ID);
		pusher.logger.Error(err.Error());		
		return err;
	}
	topic := pusher.topics[svcID.Service]
	if topic == "" {
		pusher.logger.Error(fmt.Sprintf("Unknown topic for event, event's serviceID :%d",event.CurrentAction().Service));		
	}

	pusher.channel <-&kafka.Message{ 
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Value: data,
	}
	return nil;
}

func (speaker *KafkaSpeaker) ConfigTopic(configs []microservice.MicroserviceListenerConfig) {
	for index,config := range configs {
		speaker.topics[index] = config.GetTopic();
	}
}