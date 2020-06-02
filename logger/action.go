package logger

import (
	"bytes"
	"encoding/json"
	"github.com/marsli9945/go-websocket/form"
	"log"
	"net/http"
	"time"
)

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
	resp, err := client.Post("http://47.95.216.127:9264/grow-analytics-log-server/log/send", "application/json;charset=utf-8;", bytes.NewBuffer(jsonStr))
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
}
