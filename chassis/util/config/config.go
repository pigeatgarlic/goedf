package config

import (
	"errors"

	"github.com/joeshaw/envdecode"
)

type SecurityConfig struct {
	ValidatorUrl string				`env:"VALIDATOR_URL",required`
}
func GetSecurityConfig()(*SecurityConfig, error) {
	var c SecurityConfig; 
	if err := envdecode.StrictDecode(&c); err != nil {
		return nil,err;
	}
	return &c,nil;
}






type GatewayConfig struct {
	HelperPort 	int 				`env:"HELPER_PORT",required`
	WatcherPort int 				`env:"GATEWAY_PORT",required`
}	
func GetGatewayConfig()(*GatewayConfig, error) {
	var c GatewayConfig; 
	if err := envdecode.StrictDecode(&c); err != nil {
		return nil,err;
	}
	return &c,nil;
}



type MQInterface interface {
	GetAck() string
	GetClientID() string
	GetProtocol() string
	GetProvider() string
	GetServer() string
}
type KafkaConfig struct {
	Server	string		`env:"KAFKA_SERVER",required`
	ClientID string		`env:"SERVICE_NAME",required`
	Acks string 		`env:"ACKS",required`
	DataProtocol string	`env:"DATA_PROTOCOL",required`
}
func (kafka *KafkaConfig) GetAck() string {
	return kafka.Acks;
}
func (kafka *KafkaConfig) GetClientID() string {
	return kafka.ClientID;
}
func (kafka *KafkaConfig) GetProtocol() string {
	return kafka.DataProtocol;
}
func (kafka *KafkaConfig) GetServer() string {
	return kafka.Server;
}
func (kafka *KafkaConfig) GetProvider() string {
	return "kafka";
}


func GetMQConfig()(MQInterface, error) {
	var c KafkaConfig; 
	if err := envdecode.StrictDecode(&c); err != nil {
		return nil,err;
	}
	return &c,nil;
}


type EtcdConfig struct {
	Server	string		`env:"ETCD_ADDRESS",required`
	Username string		`env:"ETCD_USER",required"`
	Password string		`env:"ETCD_PASSWORD",required"`
}

func (etcd *EtcdConfig) GetServer() string {
	return etcd.Server;
}
func (etcd *EtcdConfig) GetPassword() string {
	return etcd.Password;
}
func (etcd *EtcdConfig) GetUser() string {
	return etcd.Username;
}
func (etcd *EtcdConfig) GetChannel() string {
	return "";
}
func (etcd *EtcdConfig) GetProvider() string {
	return "etcd";
}


type RedisConfig struct {
	Server	string		`env:"REDIS_ADDRESS",required`
	Channel string		`env:"REDIS_CHANNEL",required`
}
func (etcd *RedisConfig) GetServer() string {
	return etcd.Server;
}
func (etcd *RedisConfig) GetUser() string {
	return etcd.Channel;
}
func (etcd *RedisConfig) GetPassword() string {
	return "redis";
}
func (etcd *RedisConfig) GetChannel() string {
	return etcd.Channel;
}
func (etcd *RedisConfig) GetProvider() string {
	return "redis";
}


type CacheInterface interface {
	GetServer() string
	GetChannel() string
	GetProvider() string
	GetUser() string
	GetPassword() string
}

func GetCacheConfig()(CacheInterface, error) {
	var c RedisConfig; 
	var e EtcdConfig; 

	envdecode.StrictDecode(&c); 
	envdecode.StrictDecode(&e); 

	if c.Channel != "" || c.Server != "" {
		return &c,nil;
	} else if e.Server != "" {
		return &e,nil;
	} else {
		return nil,errors.New("fail to decode env");
	}
}




type PubsubConfig struct {
	Server	string		`env:"REDIS_BROADCAST_ADDRESS",required`
	Channel string		`env:"REDIS_CHANNEL_CHANNEL",required`
}

func GetPubsubConfig()(*PubsubConfig, error) {
	var c PubsubConfig; 
	if err := envdecode.StrictDecode(&c); err != nil {
		return nil,err;
	}
	return &c,nil;
}




type ESLogConfig struct {
	ESurl        string `env:"ELASTICSEARCH,required"`
	WarningIndex string `env:"WARNING_INDEX,required"`
	ErrorIndex   string `env:"ERROR_INDEX,required"`
	InforIndex   string `env:"INFOR_INDEX,required"`
	HostName     string `env:"HOSTNAME,required"`

	StdLog    string 	`env:"STDENABLE,required"`
	Namespace string
}

func GetESlogConfig() (*ESLogConfig,error) {
	var ret ESLogConfig;
	if err := envdecode.Decode(&ret); err != nil {
		return nil,err
	}
	return &ret,nil;
}


type EventQueryConfig struct {
	Provider string
	Topic string
}

type EventPusherConfig struct {
	Provider string
	Topic string
}