package etcd

import (
	"testing"

	"github.com/pigeatgarlic/ideacrawler/microservice/chassis/util/config"
	eslogger "github.com/pigeatgarlic/ideacrawler/microservice/chassis/util/logger/es"
)

func TestEtcdConn(t *testing.T) {
	var etcdconfig = config.EtcdConfig{
		Server: "etcd-1653542552-headless:2379",
		Username: "root",
		Password: "wfJEbBCA11",
	}

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
	cache, err := InitEtcdCache(&etcdconfig, logger)
	if err != nil {
		t.Errorf(err.Error())
	}

	cache.SetKeyValue("testa", "testb")
	result, err := cache.GetKeyValue("testa")
	if err != nil {
		t.Errorf(err.Error())
	}

	if result != "testb" {
		t.Errorf("wrong value")
	}
}
