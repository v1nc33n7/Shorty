package main

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

const CacheDuration = 30

type Redis struct {
	db   *redis.Client
	pipe chan KeyValue
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
	}), make(chan KeyValue)}

	_, err := rdc.db.Ping(ctx).Result()
	if err != nil {
		return err
	}

	go rdc.Run()

	return nil
}

func (r *Redis) Run() {
	defer r.db.Conn().Close()
	defer close(r.pipe)

	for {
		msg := <-r.pipe
		err := r.CacheUrl(msg)
		if err != nil {
			log.Fatalf("%v", err)
			break
		}
	}
}

func (r *Redis) CacheUrl(msg KeyValue) error {
	err := r.db.Set(ctx, msg.Key, msg.Value, CacheDuration*time.Second).Err()

	return err
}

func (r *Redis) FindCacheUrl(msg *KeyValue) error {
	val, err := r.db.Get(ctx, msg.Key).Result()
	msg.Value = val

	return err
}
