package broker

import (
	"fmt"
	"github.com/go-redis/redis"
	"encoding/json"
	"houston/socket"
	"strconv"
	"github.com/gobwas/ws/wsutil"
	"github.com/gobwas/ws"
	"log"
)

type Message struct {
	Content string `json:"content"`
	Receivers []int `json:"receivers"`
}

type Client struct {
	*redis.Client
}

func RedisClient() Client {
	println("redis")
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
	return Client{client}
}

func (c *Client) Start() {

	pubsub := c.Subscribe("messages")
	defer pubsub.Close()
	for {
		msg, err := pubsub.ReceiveMessage()
		if err != nil {
			panic(err)
		}

		res := new(Message)
		if err := json.Unmarshal([]byte(msg.Payload), res); err != nil {
			println(err)
		}

		for _, i := range res.Receivers {
			conn := socket.ConnMap[strconv.Itoa(i)]
			if conn != nil {
				err = wsutil.WriteServerMessage(conn, ws.OpText, []byte(res.Content))
				log.Fatal(err)
			}
		}
	}

}

