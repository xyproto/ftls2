package main

import (
	"bytes"

	"github.com/garyburd/redigo/redis"
)

// Functions for dealing with a short list of string values in Redis

type RedisList struct {
	c  redis.Conn
	id string
}

func NewRedisList(c redis.Conn, id string) *RedisList {
	var rl RedisList
	rl.c = c
	rl.id = id
	return &rl
}

// Connect to the local instance of Redis at port 6379
func NewRedisConnection() (redis.Conn, error) {
	return redis.Dial("tcp", ":6379")
}

func (rl *RedisList) Store(value string) error {
	_, err := rl.c.Do("RPUSH", rl.id, value)
	return err
}

func bytes2string(b []uint8) string {
	return bytes.NewBuffer(b).String()
}

func getString(bi []interface{}, i int) string {
	return bytes2string(bi[i].([]uint8))
}

func (rl *RedisList) GetAll() ([]string, error) {
	result, err := redis.Values(rl.c.Do("LRANGE", rl.id, "0", "-1"))
	strs := make([]string, len(result))
	for i := 0; i < len(result); i++ {
		strs[i] = getString(result, i)
	}
	return strs, err
}

func (rl *RedisList) GetLast() (string, error) {
	result, err := redis.Values(rl.c.Do("LRANGE", rl.id, "-1", "-1"))
	if len(result) == 1 {
		return getString(result, 0), err
	}
	return "", err
}

func (rl *RedisList) DelAll() error {
	_, err := rl.c.Do("DEL", rl.id)
	return err
}
