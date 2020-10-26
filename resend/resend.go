package resend

import (
	"github.com/marsli9945/go-websocket/impl"
	"github.com/marsli9945/go-websocket/logger"
	"log"
	"time"
)

type resendStuck struct {
	exper int64
	List  []string
}

var resendList = map[string]*resendStuck{}

func InitFlush() {
	for {
		time.Sleep(time.Second * 5)
		log.Println("resendList", resendList)
		for k, v := range resendList {
			if v.exper+15 <= time.Now().Unix() {
				delete(resendList, k)
			}
		}
	}
}

func Add(name string, requestId string) {
	log.Println("resend add", name, requestId)
	if _, ok := resendList[name]; !ok {
		resendList[name] = &resendStuck{
			time.Now().Unix(),
			[]string{
				requestId,
			},
		}
	} else {
		resendList[name].exper = time.Now().Unix()
		resendList[name].List = append(resendList[name].List, requestId)
	}
}

func Consume(conn *impl.Connection, name string) {
	log.Println("resend consume", name)
	if resend, ok := resendList[name]; ok {
		go logger.ResendList(conn, resend.List)
		delete(resendList, name)
	}
}

func Get() map[string]*resendStuck {
	return resendList
}
