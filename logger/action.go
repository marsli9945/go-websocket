package logger

import (
	"bytes"
	"encoding/json"
	"github.com/marsli9945/go-websocket/form"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

var gapi_host = os.Getenv("GAPI_HOST") + "/api/ga/v1/grow-analytics-log-server/log/send"

func Push(event string, param form.SendForm) {
	lib := NewLib(param.Service_version)
	properties := NewProperties(param.Project_id, param.Model_name, param.Request_id)
	logContent := NewLogContent(event, param.Device_id, param.User_id, *properties, *lib)

	jsonStr, err := json.Marshal(logContent)
	if err != nil {
		log.Println(err)
	}

	// 超时时间：5秒
	client := &http.Client{Timeout: 5 * time.Second}
	log.Println("gapi_host:{}", gapi_host)
	resp, err := client.Post(gapi_host, "application/json;charset=utf-8;", bytes.NewBuffer(jsonStr))
	defer resp.Body.Close()
	if err != nil {
		log.Println(err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	log.Println("body:{}", body)
}
