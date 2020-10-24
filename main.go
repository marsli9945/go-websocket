package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/marsli9945/go-websocket/form"
	"github.com/marsli9945/go-websocket/impl"
	"github.com/marsli9945/go-websocket/logger"
	"github.com/marsli9945/go-websocket/resend"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

var (
	upgrader = websocket.Upgrader{
		// 读取存储空间大小
		ReadBufferSize: 1024,
		// 写入存储空间大小
		WriteBufferSize: 1024,
		// 允许跨域
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	// 在线用户和链接凭据
	userList = map[string]*impl.Connection{}
)

// http返回参数
type result struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	var (
		wsConn *websocket.Conn
		err    error
		// data []byte
		conn *impl.Connection
		data []byte
	)
	// 完成http应答，在httpheader中放下如下参数
	if wsConn, err = upgrader.Upgrade(w, r, nil); err != nil {
		return // 获取连接失败直接返回
	}

	if conn, err = impl.InitConnection(wsConn); err != nil {
		go logger.Push("socket_server_connect_failed", form.SendForm{})
		goto ERR
	}

	// 连接成功
	go logger.Push("socket_server_connect_success", form.SendForm{})

	//go func() {
	//	var (
	//		err error
	//	)
	//	for {
	//		// 每隔一秒发送一次心跳
	//		if err = conn.WriteMessage([]byte(`{"data": []}`)); err != nil {
	//			return
	//		}
	//		time.Sleep(1 * time.Second)
	//	}
	//}()

	for {
		if data, err = conn.ReadMessage(); err != nil {
			goto ERR
		}

		var param form.SendForm

		if err := json.Unmarshal(data, &param); err != nil {
			log.Println(err)
		}

		if param.Socket_method == "login" {
			log.Println("+++++++++++注册上线：" + param.Name)
			conn.Name = param.Name
			userList[param.Name] = conn
			resend.Consume(conn, param.Name)
		}

		if err = conn.WriteMessage(data); err != nil {
			goto ERR
		}
	}

ERR:
	// 关闭当前连接
	conn.Close()
}

func main() {

	// 开启重发队列定时清理
	go resend.InitFlush()

	// 当有请求访问ws时，执行此回调方法
	http.HandleFunc("/websocket", wsHandler)

	// 消息推动的对外http接口
	http.HandleFunc("/websocket/push", func(writer http.ResponseWriter, request *http.Request) {
		body, err := ioutil.ReadAll(request.Body)
		if err != nil {
			log.Println(err)
		}

		log.Println("param:{}" + string(body))
		var r []byte

		var param form.SendForm

		err = json.Unmarshal(body, &param)
		if err != nil {
			log.Println(err)
		}

		go logger.Push("socket_server_push_data_start", param)

		log.Println(param.Name + "+++++++开始推送")

		if param.Name == "" {
			go logger.Push("socket_server_push_data_failed", param)
			r, _ = json.Marshal(&result{401, "请使用name参数指定接收人", nil})
		} else {
			v, ok := userList[param.Name]

			if !ok {
				go logger.Push("socket_server_push_data_failed", param)
				log.Println(param.Name + "------未上线")
				r, _ = json.Marshal(&result{401, param.Name + "已断开链接", nil})
				_, _ = writer.Write(r)
				return
			}

			again := 1
			isOnline := true
			for v.IsClosed {
				if again > 30 {
					isOnline = false
					break
				}
				log.Printf("%s重试第%d次", param.Name, again)
				time.Sleep(1 * time.Second)
				again++
			}

			if isOnline {
				err = userList[param.Name].WriteMessage(body)
				if err != nil {
					log.Println(err)
					go logger.Push("socket_server_push_data_failed", param)
					resend.Add(param.Name, param.Request_id)

					log.Println(param.Name + "+++++++发送失败")
					r, _ = json.Marshal(&result{10, param.Name + "推送失败", nil})
				}
				go logger.Push("socket_server_push_data_success", param)
				log.Println(param.Name + "+++++++发送成功")
				r, _ = json.Marshal(&result{0, param.Name + "推送成功", nil})
			} else {
				go logger.Push("socket_server_push_data_failed", param)
				resend.Add(param.Name, param.Request_id)
				delete(userList, param.Name) // 清理断开的连接
				log.Println(param.Name + "------未上线")
				r, _ = json.Marshal(&result{401, param.Name + "已断开链接", nil})
			}
		}

		_, _ = writer.Write(r)
	})

	// 获取在线的用户列表
	http.HandleFunc("/websocket/list", func(writer http.ResponseWriter, request *http.Request) {
		var users []string
		for k, v := range userList {
			if !v.IsClosed {
				users = append(users, k)
			} else {
				delete(userList, k) //清理断开的连接
			}
		}
		list, _ := json.Marshal(users)
		_, _ = writer.Write(list)
	})

	http.HandleFunc("/websocket/hasUser", func(writer http.ResponseWriter, request *http.Request) {
		err := request.ParseForm()
		if err != nil {
			log.Println(err)
		}

		var r []byte
		user := request.Form.Get("name")
		if user == "" {
			r, _ = json.Marshal(&result{401, "请使用name参数指定查询用户", nil})
		} else {
			connection := userList[user]
			if connection == nil || connection.IsClosed {
				r, _ = json.Marshal(&result{0, "操作成功", false})
			} else {
				r, _ = json.Marshal(&result{0, "操作成功", true})
			}
		}
		writer.Header().Add("Access-Control-Allow-Origin", "*")
		writer.Header().Add("Access-Control-Allow-Methods", "*")
		writer.Header().Add("Access-Control-Allow-Headers", "*")
		writer.Write(r)
	})

	http.HandleFunc("/websocket/del", func(writer http.ResponseWriter, request *http.Request) {
		for k, _ := range userList {
			delete(userList, k)
		}
		r, _ := json.Marshal(&result{0, "操作成功", true})
		writer.Write(r)
	})

	http.HandleFunc("/websocket/test", func(writer http.ResponseWriter, request *http.Request) {
		err := request.ParseForm()
		if err != nil {
			log.Println(err)
		}
		key := request.Form.Get("key")
		name := request.Form.Get("name")

		strArr := strings.FieldsFunc(key, func(r rune) bool {
			if r == ',' {
				return true
			} else {
				return false
			}
		})

		log.Println(strArr)
		log.Println(name)

		logger.ResendList(userList[name], strArr)
	})

	// 渲染html文件进行测试
	http.HandleFunc("/websocket/ws", func(writer http.ResponseWriter, request *http.Request) {
		http.ServeFile(writer, request, "html/index.html")
	})
	http.HandleFunc("/websocket/wss", func(writer http.ResponseWriter, request *http.Request) {
		http.ServeFile(writer, request, "html/wss.html")
	})

	// 监听127.0.0.1:7777
	err := http.ListenAndServe("0.0.0.0:7777", nil)
	if err != nil {
		log.Fatal("ListenAndServe", err.Error())
	}
}
