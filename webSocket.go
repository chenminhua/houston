package main

import (
	"github.com/gobwas/ws"
	"net/http"
	"github.com/gobwas/ws/wsutil"
	"net"
	"houston/broker"
	"log"
	"reflect"
	"encoding/json"
	"strconv"
)

// https://medium.freecodecamp.org/million-websockets-and-go-cc58418460bb
// var ws = new WebSocket("ws://localhost:8087")
// ws.addEventListener("message", function(event) {console.log(event.data)})
// ws.send("hello")

// 如何去标示一个websocket连接呢？？
// connect的时候去http header里面拿token，解开token后把userId存下来

var connPool = make(map[string]net.Conn)

func main() {
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

	http.ListenAndServe(":8087", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 在UpGrade之前完成鉴权
		authToken := r.URL.Query().Get("auth")
		conn, _, _, err := ws.UpgradeHTTP(r, w, nil)
		connPool[authToken] = conn
		HandleError(err)

		// 给每个长连接创建一个goroutine (仅仅为了监视websocket的状态？？)
		go func() {
			defer cleanUpConn(conn, authToken)
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
		}()
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



