package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	*redis.Client
}

var (
	ctx = context.Background()
	rdc *Redis
)

func ConnRedis(addr string) error {
	rdc = &Redis{redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})}

	_, err := rdc.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("couldn't connect to redis server")
	}

	err = rdc.enableNotify()
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	return nil
}

func (r *Redis) enableNotify() error {
	err := r.ConfigSet(ctx, "notify-keyspace-events", "KEA")
	if err != nil {
		//return fmt.Errorf("couldn't config notify-keyspace-events in redis")
	}

	pubsub := r.PSubscribe(ctx, "__keyevent@0__:expired")
	go func(redis.PubSub) {
		for {
			msg, err := pubsub.ReceiveMessage(ctx)
			if err != nil {
				log.Printf("Redis [PubSub]: Error message: %v", err)
			}
			log.Printf("Redis [PubSub]: Recieved: %v", msg)
		}
	}(*pubsub)

	return nil
}

func (r *Redis) AddUrl(newUrl string, orginal string) error {
	err := r.Set(ctx, newUrl, orginal, 10*time.Second).Err()
	if err != nil {
		log.Printf("Redis: Couldn't add new key %s, value %s", newUrl, orginal)
		return err
	}

	return nil
}

func (r *Redis) FindUrl(url string) (string, error) {
	val, err := r.Get(ctx, url).Result()
	if err != nil {
		log.Printf("Redis: Couldn't read key %s", url)
		return "", err
	}

	return val, nil
}
