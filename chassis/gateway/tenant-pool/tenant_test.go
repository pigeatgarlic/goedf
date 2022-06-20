package tenantpool

import (
	"testing"

	"github.com/pigeatgarlic/goedf/chassis/gateway/tenant-pool/tenant"
	"github.com/pigeatgarlic/goedf/chassis/util/config"
	eslogger "github.com/pigeatgarlic/goedf/chassis/util/logger/es"
	"github.com/pigeatgarlic/goedf/models/request-response/response"
	"github.com/pigeatgarlic/goedf/models/user"
)

func TestSingleTenantPool(t *testing.T) {
	logconf := &config.ESLogConfig{
		ESurl: "http://elasticsearch-1653540188.elasticstack.svc.cluster.local:9200",
		WarningIndex: "testwarning",
		InforIndex: "testinfor",
		ErrorIndex: "testerror",
		HostName: "unittest",
		StdLog: "true",
		Namespace: "unittest",
	}

	logger,err :=	eslogger.InitLogger(logconf);
	if err != nil {
		t.Error(err)
		return;
	}


	pool,err := InitTenantPool(&config.PubsubConfig{
		Server: "redis-1653548263:6379",
		Channel: "test",
	},logger);
	if err != nil {
		t.Error(err)
		return;
	}

	new := tenant.NewTenant(0,&user.User{
		ID: 0,
		UserName: "test",
	})


	pool.NewTenant(new);
	pool.SendResponse(&response.UserResponse{
		ID: 0,
		SessionID: 0,
		Error: "",
		Data: map[string]string {
			"data": "test",
		},
	})

	logger.Infor("Waiting for response")
	result := new.ListenonResponse();
	if result.Data["data"] != "test" {
		t.Error("unmatch result")
		return;
	}
	pool.KillTenant(0);
}


func TestMultipleTenantPool(t *testing.T) {
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
	pool1,err := InitTenantPool(&config.PubsubConfig{
		Server: "redis-1653548263:6379",
		Channel: "test",
	},logger);
	if err != nil {
		t.Error(err)
		return;
	}

	pool2,err := InitTenantPool(&config.PubsubConfig{
		Server: "redis-1653548263:6379",
		Channel: "test",
	},logger);
	if err != nil {
		t.Error(err)
		return;
	}

	new := tenant.NewTenant(0,&user.User{
		ID: 0,
		UserName: "test",
	})

	pool1.NewTenant(new);

	pool2.SendResponse(&response.UserResponse{
		ID: 0,
		SessionID: 0,
		Error: "",
		Data: map[string]string {
			"data": "test",
		},
	})

	logger.Infor("Waiting for response")
	result := new.ListenonResponse();
	if result.Data["data"] != "test" {
		t.Error("unmatch result")
	}
}