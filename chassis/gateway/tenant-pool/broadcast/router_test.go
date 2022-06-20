package broadcast

import (
	"testing"

	redispubsub "github.com/pigeatgarlic/goedf/chassis/gateway/tenant-pool/broadcast/redis"
	"github.com/pigeatgarlic/goedf/chassis/util/config"
	eslogger "github.com/pigeatgarlic/goedf/chassis/util/logger/es"
)

func TestRedisRouter(t *testing.T) {
	logconf := &config.ESLogConfig{
		ESurl: "http://elasticsearch-1653540188.elasticstack.svc.cluster.local:9200",
		WarningIndex: "testwarning",
		InforIndex: "testinfor",
		ErrorIndex: "testerror",
		HostName: "unittest",
		StdLog: "true",
		Namespace: "unittest",
	}

	logger,_ :=	eslogger.InitLogger(logconf);

	pubsub,err := redispubsub.InitRedisPubSub(&config.PubsubConfig{
		Server: "redis-1653548263:6379",
		Channel: "test",
	},logger);
	if err != nil  {
		t.Error(err)
		return;
	}

	testmsg := []byte("test");
	pubsub.Publish([]byte(testmsg))
	ret := pubsub.Subscribe()
	if string(ret) == string(testmsg) {
		return;
	}
	t.Error("Unmatch result");
}