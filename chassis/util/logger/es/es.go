package eslogger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"time"

	"github.com/pigeatgarlic/goedf/chassis/util/config"
)

type inforMessage struct {
	Message   string
	Time      time.Time
	Server    string
	Namespace string
}

type warningMessage struct {
	Message   string
	Time      time.Time
	Server    string
	Namespace string
}

type errorMessage struct {
	StackTrace string
	Message    string
	Time       time.Time
	Server     string
	Namespace  string
}

type EsLogger struct {
	config     *config.ESLogConfig
	httpClient http.Client

	elasticLog bool
}

// init new index on elasticsearch by create a PUT request
func (logger *EsLogger) initIndex(index string) {
	resp, err := logger.httpClient.Get(logger.config.ESurl + "/" + index)
	if err != nil {
		fmt.Println("Fail to create new elasticsearch index")
	} else {
		if resp.StatusCode == 404 {
			req, _ := http.NewRequest(http.MethodPut, logger.config.ESurl+"/"+index, bytes.NewBufferString(""))
			resp, err := logger.httpClient.Do(req)
			if err != nil || resp.StatusCode != 200 {
				fmt.Println("Fail to create new elasticsearch index")
			}
		}
	}
}

// Set logger as global variable to make it an singleton
func InitLogger(config *config.ESLogConfig) (*EsLogger, error) {
	var logger EsLogger

	logger.elasticLog = true

	logger.httpClient = http.Client{Timeout: 1 * time.Second}
	logger.config = config

	resp, err := http.Get(logger.config.ESurl)
	if err != nil || resp.StatusCode != 200 {
		log.Println("Fail to connect to elasticsearch due to ", err, ", pushing log to stdout")
		logger.elasticLog = false
	} else {
		logger.initIndex(logger.config.WarningIndex)
		logger.initIndex(logger.config.InforIndex)
		logger.initIndex(logger.config.ErrorIndex)
	}

	file, namespace_err := os.ReadFile("/run/secrets/kubernetes.io/serviceaccount/namespace")
	if namespace_err != nil {
		logger.config.Namespace = "unknown"
	} else {
		logger.config.Namespace = string(file)
	}

	return &logger, nil
}

func FreeLogger() {

}

func (logger *EsLogger) pushLog(data []byte, index string) {
	// if !logger.elasticLog {
	// 	return
	// }

	body := bytes.NewBuffer(data)
	resp, err := logger.httpClient.Post(fmt.Sprintf("%s/%s/_doc", logger.config.ESurl, index), "application/json", body)
	if err != nil {
		fmt.Println("Failed to push log to elasicsearch, error: ", err)
	}

	if resp.StatusCode == 201 {
		// fmt.Println("success");
	} else {
		// res,_ := io.ReadAll(resp.Body)
		// fmt.Println("%d:  %s", resp.StatusCode, string(res));
	}
}

func (logger *EsLogger) Fatal(format string) {
	var message errorMessage
	message.Message = format
	message.Time = time.Now()
	message.StackTrace = string(debug.Stack())
	message.Namespace = logger.config.Namespace
	message.Server = logger.config.HostName

	fmt.Printf("%s  %s  %s\n", message.Time, message.Message, message.StackTrace)

	data, _ := json.Marshal(message)
	go logger.pushLog(data, logger.config.ErrorIndex)
	panic(message.Message)
}

func (logger *EsLogger) Debug(format string) {
	var message inforMessage
	message.Message = format
	message.Time = time.Now()
	message.Namespace = logger.config.Namespace
	message.Server = logger.config.HostName

	fmt.Printf("%v  %s\n", message.Time, message.Message)

	data, _ := json.Marshal(message)
	go logger.pushLog(data, logger.config.InforIndex)
}

func (logger *EsLogger) Error(format string) {
	var message errorMessage
	message.Message = format
	message.Time = time.Now()
	message.StackTrace = string(debug.Stack())
	message.Namespace = logger.config.Namespace
	message.Server = logger.config.HostName

	fmt.Printf("%s  %s  %s\n", message.Time, message.Message, message.StackTrace)

	data, _ := json.Marshal(message)
	go logger.pushLog(data, logger.config.ErrorIndex)
}

func (logger *EsLogger) Warning(format string) {
	var message warningMessage
	message.Message = format
	message.Time = time.Now()
	message.Namespace = logger.config.Namespace
	message.Server = logger.config.HostName

	fmt.Printf("%v  %s\n", message.Time, message.Message)

	data, _ := json.Marshal(message)
	go logger.pushLog(data, logger.config.WarningIndex)
}

func (logger *EsLogger) Infor(format string) {
	var message inforMessage
	message.Message = format
	message.Time = time.Now()
	message.Namespace = logger.config.Namespace
	message.Server = logger.config.HostName

	fmt.Printf("%v  %s\n", message.Time, message.Message)

	data, _ := json.Marshal(message)
	go logger.pushLog(data, logger.config.InforIndex)
}

func (logger *EsLogger) Find(index string, format string, start *string, end *string) []byte {
	var body []byte

	if start == nil && end == nil {
		query_string := fmt.Sprintf(
			`{
			"query": {
				"term" : {
					"Message" : "%s"
				}
			}
		}`, format)

		body = []byte(query_string)
	}

	logger.Debug(fmt.Sprintf("Finding %s keyword on index name %s", format, index))
	resp, err := logger.httpClient.Post(fmt.Sprintf("%s/%s/_search", logger.config.ESurl, index), "application/json", bytes.NewBuffer(body))
	if err != nil {
		logger.Error("Post to elasicsearch, error: " + err.Error())
		return nil
	}

	resp_data, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err.Error())
		return nil
	}

	return resp_data
}
