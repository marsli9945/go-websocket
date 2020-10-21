package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func main() {
	gapi_host := "https://gapics.touch4.me/api/ga/v1/grow-analytics-log-server/log/send"
	jsonStr := `{"Client_id":"H5_5.0_tuyoo.tuyoo.0-hall20435.tuyoo.GA","Device_id":"lilei@tuyoogame.com","User_id":"91","Event":"socket_server_push_data_start","Event_time":1603283155,"Project_id":"20435","type":"track","Properties":{"Proj_project_id":"20249","Proj_model_version":"0.1.0","Proj_service_name":"websocket","Proj_model_name":"report","Proj_request_id":"1603283151179-3ea14f8aa5c4488c989e1c165b7c9105","Proj_cost_time":4124},"Lib":{"Lib_service_version":"","Lib_language":"GO"}}`
	byte := []byte(jsonStr)

	// 超时时间：5秒
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Post(gapi_host, "application/json;charset=utf-8;", bytes.NewBuffer(byte))
	if err != nil {
		log.Println(err)
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		log.Println("body:{}", string(body))
	}
}
