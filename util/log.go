package util

import "time"

type LogContent struct {
	Client_id  string
	Request_id string
	Device_id  string
	User_id    string
	Event      string
	Event_time string
	Project_id string
	Event_type string `json:"type"`
	Properties Properties
	Pib        Lib
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
	Proj_search_id     int64
	Proj_cost_time     int64
}

func NewProperties(proj_project_id string, proj_model_name string, proj_search_id int64) *Properties {
	return &Properties{
		Proj_project_id:    proj_project_id,
		Proj_model_version: "0.1.0",
		Proj_service_name:  "websocket",
		Proj_model_name:    proj_model_name,
		Proj_search_id:     proj_search_id,
		Proj_cost_time:     time.Now().Unix() - proj_search_id,
	}
}
