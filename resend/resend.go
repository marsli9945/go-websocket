package resend

import (
	"github.com/marsli9945/go-websocket/form"
	"github.com/marsli9945/go-websocket/impl"
	"github.com/marsli9945/go-websocket/logger"
	"log"
	"time"
)

type resendStuck struct {
	Device_id string
	User_id   string
	exper     int64
	List      []string
}

var resendList = map[string]*resendStuck{}

func InitFlush() {
	for {
		time.Sleep(time.Second * 5)
		log.Println("resendList", resendList)
		for k, v := range resendList {
			if v.exper+60 <= time.Now().Unix() {
				for _, rv := range v.List {
					go logger.Push("socket_server_clean_resend", form.SendForm{
						Request_id: rv,
						Device_id:  v.Device_id,
						User_id:    v.User_id,
					})
				}
				delete(resendList, k)
			}
		}
	}
}

func Add(form form.SendForm) {
	log.Println("resend add", form.Name, form.Request_id)
	go logger.Push("socket_server_resend_add", form)
	if _, ok := resendList[form.Name]; !ok {
		resendList[form.Name] = &resendStuck{
			form.Device_id,
			form.User_id,
			time.Now().Unix(),
			[]string{
				form.Request_id,
			},
		}
	} else {
		resendList[form.Name].exper = time.Now().Unix()
		resendList[form.Name].List = append(resendList[form.Name].List, form.Request_id)
	}
}

func Consume(conn *impl.Connection, name string) {
	log.Println("resend consume", name)
	if resend, ok := resendList[name]; ok {
		go ResendList(name, conn, resend)
		delete(resendList, name)
	}
}

func Get() map[string]*resendStuck {
	return resendList
}
