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
		return fmt.Errorf("%v", err)
	}

	return nil
}

func (r *Redis) CacheUrl(msg KeyValue) error {
	err := r.Set(ctx, msg.Key, msg.Value, 30*time.Second).Err()
	if err != nil {
		log.Printf("Redis: Couldn't add new key %s, value %s", msg.Key, msg.Value)
		return err
	}

	return nil
}

func (r *Redis) FindCacheUrl(url string) (string, error) {
	val, err := r.Get(ctx, url).Result()
	if err != nil {
		log.Printf("Redis: Couldn't read key %s", url)
		return "", err
	}

	return val, nil
}
