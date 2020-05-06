package form

type SendForm struct {
	Name            string      `json:"name"`
	Data            interface{} `json:"data"`
	Project_id      int16       `json:"project_id"`
	Request_id      int64       `json:"request_id"`
	Model_name      string      `json:"model_name"`
	Service_version string      `json:"service_version"`
	Device_id       string      `json:"device_id"`
	User_id         int16       `json:"user_id"`
}
