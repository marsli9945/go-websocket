package logger

import (
	"bytes"
	"encoding/json"
	"github.com/marsli9945/go-websocket/form"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

//var gapi_host = GapiHost + "/v1/grow-analytics-log-server/log/record"
var gapi_host = GapiHost + "/api/ga/v1/grow-analytics-log-server/log/record"

type SendParam struct {
	Data *LogContent `json:"data"`
}

var logList = []*LogContent{}

func Push(event string, param form.SendForm) {
	lib := NewLib(param.Service_version)
	properties := NewProperties(param.Project_id, param.Model_name, param.Request_id)
	logContent := NewLogContent(event, param.Device_id, param.User_id, *properties, *lib)
	logList = append(logList, logContent)
	log.Printf("logContent", logContent)
	if len(logList) < 20 {
		return
	}

	jsonStr, err := json.Marshal(logList)
	if err != nil {
		log.Println(err)
		return
	}
	logList = []*LogContent{}

	// 超时时间：2秒
	client := &http.Client{Timeout: 2 * time.Second}
	log.Println("jsonStr:{}", string(jsonStr))
	resp, err := client.Post(gapi_host, "application/json;charset=utf-8;", bytes.NewBuffer(jsonStr))
	if err != nil {
		log.Println(err)
	} else {
		defer resp.Body.Close()
		all, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("send error:", err)
		}
		log.Println("body:{}", string(all))
	}
}
