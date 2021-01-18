package tool

import (
	"log"
	"time"

	redigo "github.com/gomodule/redigo/redis"
)

var pool *redigo.Pool

func NewRedisPool() *redigo.Pool {
	cfg := GetConfig()
	return &redigo.Pool{
		MaxIdle:     cfg.Redis.MaxIdle,
		MaxActive:   cfg.Redis.MaxActive,
		IdleTimeout: time.Duration(cfg.Redis.IdleTimeout) * time.Second,
		Dial: func() (redigo.Conn, error) {
			c, err := redigo.Dial("tcp", cfg.Redis.Addr+":"+cfg.Redis.Port)
			if err != nil {
				log.Fatal(err.Error())
				return nil, err
			}
			return c, nil
		},
	}
}
