package cache

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

func New(addr, password string) * redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: 		addr, 
		Password: 	password,
		DB:			0,
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("cache: failed to connect to redis at %s: %v", addr, err)
	}

	log.Println("cache: connected to redis")
	return rdb
}