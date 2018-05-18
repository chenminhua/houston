package broker

import (
	"fmt"
	"github.com/go-redis/redis"
)

type Message struct {
	Content string `json:"content"`
	Receivers []int `json:"receivers"`
}

func RedisClient() *redis.Client {
	println("redis")
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
	return client

	//go func() {
	//	pubsub := client.Subscribe("messages")
	//	defer pubsub.Close()
	//	for {
	//		msg, err := pubsub.ReceiveMessage()
	//		if err != nil {
	//			panic(err)
	//		}
	//
	//		res := new(Message)
	//		if err := json.Unmarshal([]byte(msg.Payload), res); err != nil {
	//			println(err)
	//		}
	//		fmt.Println(res)
	//	}
	//}()

}

