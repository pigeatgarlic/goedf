package registrator

import (
	"fmt"
	"testing"

	"github.com/pigeatgarlic/ideacrawler/microservice/chassis/gateway/module/registrator/etcd"
	"github.com/pigeatgarlic/ideacrawler/microservice/chassis/util/config"
	eslogger "github.com/pigeatgarlic/ideacrawler/microservice/chassis/util/logger/es"
	"github.com/pigeatgarlic/ideacrawler/microservice/models/microservice"
	"github.com/pigeatgarlic/ideacrawler/microservice/models/user"
)

// TODO add result validation
func TestEtcd(t *testing.T) {

	logger_conf := config.ESLogConfig{
		ESurl: "http://elasticsearch-1653540188.elasticstack.svc.cluster.local:9200",
		WarningIndex: "testwarning",
		InforIndex: "testinfor",
		ErrorIndex: "testerror",
		HostName: "unittest",
		StdLog: "true",
		Namespace: "unittest",
	}
	conf := config.EtcdConfig{
		Server: "etcd-1653542552-headless:2379",
		Username: "root",
		Password: "wfJEbBCA11",
	}


	logger,_ := eslogger.InitLogger(&logger_conf);
	cache,err := etcd.InitEtcdCache(&conf,logger)
	if err != nil {
		t.Error(err)
		return;
	}

	err = cache.RegisterFeature(&microservice.Feature{
		ID: 0,
		Name: "test",
		Tags: map[string]string {
			"purpose": "test",
		},

		Authority: "test",
		EndpointIDs: []int{0},
		Allowed: []user.Role{{
			ID: 0,
			Name: "tester",
		}},
	})
	if err != nil {
		t.Error(err)
		return;
	}
	err = cache.RegisterMicroservice(&microservice.MicroService{
		ID: 0,
		Name: "test",
		Tags: map[string]string {
			"purpose": "test",
		},

		Endpoints: []microservice.Endpoint{{
			ID: 0,
			Name: "test",
			InstructionSet: []microservice.Instruction{},
			Order: 0,
			MicroserviceID: 0,
		}},
	})
	if err != nil {
		t.Error(err)
		return;
	}

	result,err := cache.EventLookup("test")
	if err != nil {
		t.Error(err)
		return;
	}
	fmt.Printf("%v\n",result);

}