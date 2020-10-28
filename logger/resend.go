package logger

import (
	"bytes"
	"encoding/json"
	"github.com/marsli9945/go-websocket/impl"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Token struct {
	token string
	exper int64
}

var token = Token{}
var err error

type LoginForm struct {
	ClientId     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
	Username     string `json:"username"`
	Password     string `json:"password"`
}

var loginParam = LoginForm{
	GapiClientId,
	GapiClientSecret,
	GapiUsername,
	GapiPassword,
}

//var loginParam = LoginForm{
//	"net5ijy",
//	"y4cZe@wmGBofIlQ",
//	"lilei@tuyoogame.com",
//	"Lilei@123",
//}

//var resendHost = "https://gapics.touch4.me/api/ga/v1/insight/data-service-a/result?key="
//var loginHost = "https://gapics.touch4.me/api/ga/v1/auth/login"

var loginHost = GapiHost + "/api/ga/v1/auth/login"
var resendHost = GapiHost + "/api/ga/v1/insight/data-service-a/result?key="

func Send(method string, url string, form interface{}, header map[string]string) (response []byte, err error) {

	client := &http.Client{Timeout: 2 * time.Second}
	body, err := json.Marshal(form)
	if err != nil {
		log.Println(err)
	}

	request, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json;charset=utf-8;")
	for k, v := range header {
		request.Header.Set(k, v)
	}
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	} else {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		return body, nil
	}
}

func login() {
	log.Println("loginParam", loginParam)
	post, err := Send("POST", loginHost, loginParam, nil)
	if err != nil {
		log.Println(err)
	}
	result := map[string]interface{}{}
	err = json.Unmarshal(post, &result)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("login result:", result)
	if r, ok := result["data"].(map[string]interface{}); ok {
		if m, ok := r["access_token"].(string); ok {
			token.token = m
			token.exper = time.Now().Unix() + 3600
		}
	}
}

func ResendList(conn *impl.Connection, list []string) {
	if token.token == "" || time.Now().Unix() > token.exper {
		login()
	}

	body := ""
	for i, v := range list {
		body += v
		if i < len(list)-1 {
			body += ","
		}
	}

	header := map[string]string{
		"Authorization": token.token,
	}

	response, err := Send("GET", resendHost+body, nil, header)
	if err != nil {
		log.Println(err)
		return
	}

	result := map[string]interface{}{}
	err = json.Unmarshal(response, &result)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("result:", result)
	if code, ok := result["code"].(float64); ok {
		if code != 0 {
			return
		}

		if data, ok := result["data"].([]interface{}); ok {
			for _, msg := range data {
				marshal, err := json.Marshal(msg)
				if err != nil {
					log.Println(err)
					continue
				}

				err = conn.WriteMessage(marshal)
				if err != nil {
					log.Println(err)
				}
				log.Println("resend successfly")
			}
		}
	}
}
