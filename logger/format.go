package logger

import (
	"log"
	"strconv"
	"strings"
	"time"
)

type LogContent struct {
	Client_id  string
	Device_id  string
	User_id    string
	Event      string
	Event_time int64
	Project_id string
	Event_type string `json:"type"`
	Properties Properties
	Lib        Lib
}

func NewLogContent(event string, device_id string, user_id string, properties Properties, lib Lib) *LogContent {
	return &LogContent{
		Client_id:  "H5_5.0_tuyoo.tuyoo.0-hall20435.tuyoo.GA",
		Device_id:  device_id,
		User_id:    user_id,
		Event:      event,
		Event_time: time.Now().Unix(),
		Project_id: "20435",
		Event_type: "track",
		Properties: properties,
		Lib:        lib,
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
	Proj_project_id    string
	Proj_model_version string
	Proj_service_name  string
	Proj_model_name    string
	Proj_request_id    string
	Proj_cost_time     int64
}

func NewProperties(proj_project_id string, proj_model_name string, proj_request_id string) *Properties {
	strArr := strings.FieldsFunc(proj_request_id, func(r rune) bool {
		if r == '-' {
			return true
		} else {
			return false
		}
	})
	t, err := strconv.ParseInt(strArr[0], 10, 64)
	if err != nil {
		log.Println(err)
		t = 0
	}
	return &Properties{
		Proj_project_id:    proj_project_id,
		Proj_model_version: "0.1.0",
		Proj_service_name:  "websocket",
		Proj_model_name:    proj_model_name,
		Proj_request_id:    proj_request_id,
		Proj_cost_time:     time.Now().Unix() - t,
	}
}
