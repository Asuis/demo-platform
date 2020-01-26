package redis

import (
	"fmt"
	"github.com/go-redis/redis"
	"testing"
)

func TestRedisClient(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "39.104.121.176:6379", // use default Addr
		Password: "",               // no password set
		DB:       0,                // use default DB
	})

	pong, err := rdb.Ping().Result()
	fmt.Println(pong, err)
}