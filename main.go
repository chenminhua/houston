package main

import (
	"github.com/gobwas/ws"
	"net/http"
	"github.com/gobwas/ws/wsutil"
	"net"
	"houston/broker"
	"log"
	"encoding/json"
	"strconv"
	"houston/netpoll"
	"fmt"
	"reflect"
	"errors"
	"houston/restapi"
)

// 如何去标示一个websocket连接呢？？
// connect的时候去http header里面拿token，解开token后把userId存下来
// 如果同一个用户（同一个userId）连接了两个socket上来怎么处理？
// https://www.jianshu.com/p/c322edca985f

var connPool = make(map[string]net.Conn)

func main() {
	server := restapi.Server{Port: ":8081"}
	go server.Start()
	client := broker.RedisClient()
	go func() {
		pubsub := client.Subscribe("messages")
		defer pubsub.Close()
		for {
			msg, err := pubsub.ReceiveMessage()
			if err != nil {
				panic(err)
			}

			res := new(broker.Message)
			if err := json.Unmarshal([]byte(msg.Payload), res); err != nil {
				println(err)
			}

			for _, i := range res.Receivers {
				conn := connPool[strconv.Itoa(i)]
				if conn != nil {
					err = wsutil.WriteServerMessage(conn, ws.OpText, []byte(res.Content))
					HandleError(err)
				}
			}
		}
	}()
	setUpServer()
}

func setUpServer() {

	poller, err := netpoll.New(nil)
	if err != nil {
		fmt.Println(err)
	}

	http.ListenAndServe(":8087", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 在UpGrade之前完成鉴权
		authToken := r.URL.Query().Get("auth")
		conn, _, _, err := ws.UpgradeHTTP(r, w, nil)
		connPool[authToken] = conn
		HandleError(err)

		fd, _:= getFileDescriptor(conn)
		desc := netpoll.NewDesc(fd, netpoll.EventRead)

		poller.Start(desc, func(ev netpoll.Event) {
			msg, op, err := wsutil.ReadClientData(conn)
			if reflect.TypeOf(err) == reflect.TypeOf(wsutil.ClosedError{}) {
				cleanUpConn(conn, authToken)
			} else {
				HandleError(err)
			}
			if op == ws.OpText {
				err = wsutil.WriteServerMessage(conn, ws.OpText, msg)
				HandleError(err)
			}
		})

		// 给每个长连接创建一个goroutine (仅仅为了监视websocket的状态？？)
		// go receiveMessage(conn, authToken)
	}))
}

func cleanUpConn(conn net.Conn, userId string) {
	delete(connPool, userId)
	conn.Close()
}

func HandleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func getFileDescriptor(conn net.Conn) (fd uintptr, err error) {
	tcpConn, ok := conn.(*net.TCPConn)
	if !ok {
		return 0, errors.New("not a TCPConn")
	}

	file, err := tcpConn.File()
	if err != nil {
		return 0, err
	}
	defer file.Close()
	return file.Fd(), nil
}

func receiveMessage(conn net.Conn, userId string) {
	defer cleanUpConn(conn, userId)
	for {
		msg, op, err := wsutil.ReadClientData(conn)
		if reflect.TypeOf(err) == reflect.TypeOf(wsutil.ClosedError{}) {
			break;
		} else {
			HandleError(err)
		}
		if op == ws.OpText {
			err = wsutil.WriteServerMessage(conn, ws.OpText, msg)
			HandleError(err)
		}
	}
}



