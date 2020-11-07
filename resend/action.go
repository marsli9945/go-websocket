package resend

import (
	"bytes"
	"encoding/json"
	"github.com/marsli9945/go-websocket/form"
	"github.com/marsli9945/go-websocket/impl"
	"github.com/marsli9945/go-websocket/logger"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var err error

var resendHost = logger.GapiHost + "/api/ga/v1/insight/data-service-a/result?key="

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

func ResendList(name string, conn *impl.Connection, list []string) {
	body := ""
	for i, v := range list {
		body += v
		if i < len(list)-1 {
			body += ","
		}
	}

	response, err := Send("GET", resendHost+body, nil, nil)
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

				log.Println("resend start.............")
				err = conn.WriteMessage(marshal)
				if rmsg, ok := msg.(map[string]interface{}); ok {
					log.Println("rmsg:", rmsg)
					if rdataList, ok := rmsg["data"].([]interface{}); ok {
						for _, rinterface := range rdataList {
							if rdata, ok := rinterface.(map[string]interface{}); ok {
								log.Println("rdata:", rdata)
								if rcontent, ok := rdata["content"].(map[string]interface{}); ok {
									log.Println("rcontent:", rcontent)
									if rid, ok := rcontent["id"].(string); ok {
										log.Println("rid:", rid)
										loggerParams := form.SendForm{
											Device_id:  name,
											Request_id: rid,
											User_id:    "10000",
										}
										if err != nil {
											log.Println(err)
											logger.Push("socket_server_resend_failed", loggerParams)
										} else {
											logger.Push("socket_server_resend_success", loggerParams)
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
}
