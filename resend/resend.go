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
		for k, v := range resendList {
			if v.exper+15 <= time.Now().Unix() {
				log.Println(111)
				delete(resendList, k)
			}
		}
	}
}

func Add(name string, requestId string) {
	if _, ok := resendList[name]; !ok {
		resendList[name] = &resendStuck{
			time.Now().Unix(),
			[]string{
				requestId,
			},
		}
	} else {
		resendList[name].List = append(resendList[name].List, requestId)
	}
}

func Consume(conn *impl.Connection, name string) {
	if resend, ok := resendList[name]; ok {
		go logger.ResendList(conn, resend.List)
		delete(resendList, name)
	}
}

func Get() map[string]*resendStuck {
	return resendList
}
