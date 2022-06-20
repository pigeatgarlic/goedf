package logger

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/pigeatgarlic/ideacrawler/microservice/chassis/util/config"
	eslogger "github.com/pigeatgarlic/ideacrawler/microservice/chassis/util/logger/es"
)

func TestESLogger(t *testing.T) {
	conf := &config.ESLogConfig{
		ESurl: "http://elasticsearch-1653540188.elasticstack.svc.cluster.local:9200",
		WarningIndex: "testwarning",
		InforIndex: "testinfor",
		ErrorIndex: "testerror",
		HostName: "unittest",
		StdLog: "true",
		Namespace: "unittest",
	}

	logger,err := eslogger.InitLogger(conf);
	if err != nil {
		t.Error(err);
		return;
	}

	now := fmt.Sprintf("%d",(time.Now()).UTC().Nanosecond())
	logger.Infor(now);

	for {
		var ret interface{}
		result := logger.Find("testinfor",now,nil,nil)
		err = json.Unmarshal(result,&ret)
		if err != nil {
			logger.Error(err.Error())
		}
		
		json_ret := ret.(map[string]interface{})

		hits := json_ret["hits"].(map[string]interface{})

		count := hits["total"].(map[string]interface{})["value"].(float64)
		array := hits["hits"].([]interface{});

		done  := array[0].(map[string]interface{})["_source"].(map[string]interface{})["Time"].(string)

		if len(array) == int(count) && count >= 1{

			logger.Infor(fmt.Sprintf("Test done at %s",done));
			return;
		} else {
			logger.Infor(fmt.Sprintf("%v",string(result)));
		}
		
		time.Sleep(time.Second);
	}
}