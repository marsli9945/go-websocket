package form

type SendForm struct {
	Name            string      `json:"name"`
	Data            interface{} `json:"data"`
	Project_id      string      `json:"project_id"`
	Request_id      string      `json:"request_id"`
	Model_name      string      `json:"model_name"`
	Service_version string      `json:"service_version"`
	Device_id       string      `json:"device_id"`
	User_id         string      `json:"user_id"`
	Socket_method   string      `json:"socket_method"`
}
