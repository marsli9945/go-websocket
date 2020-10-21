package logger

import (
	"log"
	"strconv"
	"strings"
	"time"
)

type LogContent struct {
	client_id  string
	device_id  string
	user_id    string
	event      string
	event_time int64
	project_id string
	event_type string `json:"type"`
	properties Properties
	lib        Lib
}

func NewLogContent(event string, device_id string, user_id string, properties Properties, lib Lib) *LogContent {
	return &LogContent{
		client_id:  "H5_5.0_tuyoo.tuyoo.0-hall20435.tuyoo.GA",
		device_id:  device_id,
		user_id:    user_id,
		event:      event,
		event_time: time.Now().Unix(),
		project_id: "20435",
		event_type: "track",
		properties: properties,
		lib:        lib,
	}
}

type Lib struct {
	Lib_service_version string
	Lib_language        string
}

func NewLib(lib_service_version string) *Lib {
	return &Lib{Lib_service_version: lib_service_version, Lib_language: "GO"}
}

type Properties struct {
	proj_project_id    string
	proj_model_version string
	proj_service_name  string
	proj_model_name    string
	proj_request_id    string
	proj_cost_time     int64
}

func NewProperties(proj_project_id string, proj_model_name string, proj_request_id string) *Properties {
	var t int64
	var err error
	t = 0
	if proj_request_id != "" {
		strArr := strings.FieldsFunc(proj_request_id, func(r rune) bool {
			if r == '-' {
				return true
			} else {
				return false
			}
		})
		t, err = strconv.ParseInt(strArr[0], 10, 64)
		if err != nil {
			log.Println(err)
			t = 0
		}
	}

	return &Properties{
		proj_project_id:    proj_project_id,
		proj_model_version: "0.1.0",
		proj_service_name:  "websocket",
		proj_model_name:    proj_model_name,
		proj_request_id:    proj_request_id,
		proj_cost_time:     time.Now().UnixNano()/1e6 - t,
	}
}
