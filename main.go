package main

import (
	"houston/broker"
	"log"
	"houston/restapi"
	"houston/socket"
)

// 如何去标示一个websocket连接呢？？
// connect的时候去http header里面拿token，解开token后把userId存下来
// 如果同一个用户（同一个userId）连接了两个socket上来怎么处理？
// https://www.jianshu.com/p/c322edca985f

func main() {
	server := restapi.Server{Port: ":8081"}
	go server.Start()
	client := broker.RedisClient()
	go client.Start()
	socket.SetUpWebSocketServer()
}

func HandleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}





