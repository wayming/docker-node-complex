package main

import (
	"os"
	"strconv"

	"github.com/go-redis/redis"
)

func fib(index int) int {
	if index < 2 {
		return index
	}
	return fib(index-1) + fib(index-2)

}
func main() {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		Password: "",
		DB:       0})

	sub := redisClient.Subscribe("insert")
	defer sub.Close()

	msg, err := sub.ReceiveMessage()
	if err != nil {
		panic(err)
	}

	index, err := strconv.Atoi(msg.Payload)
	if err != nil {
		panic(err)
	}

	cmd := redisClient.HSet("values", msg.Payload, fib(index))
	if cmd.Err() != nil {
		panic(cmd.Err())
	}

}
