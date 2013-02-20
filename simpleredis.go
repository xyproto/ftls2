package main

import (
	"bytes"

	"github.com/garyburd/redigo/redis"
)

// Functions for dealing with a short list of string values in Redis

type RedisDatastructure struct {
	c  redis.Conn
	id string
}

type RedisList RedisDatastructure
type RedisKeyValue RedisDatastructure
type RedisHashMap RedisDatastructure
type RedisSet RedisDatastructure

func NewRedisList(c redis.Conn, id string) *RedisList {
	return &RedisList{c, id}
}

func NewRedisKeyValue(c redis.Conn, id string) *RedisKeyValue {
	return &RedisKeyValue{c, id}
}

func NewRedisHashMap(c redis.Conn, id string) *RedisHashMap {
	return &RedisHashMap{c, id}
}

func NewRedisSet(c redis.Conn, id string) *RedisSet {
	return &RedisSet{c, id}
}

// Connect to the local instance of Redis at port 6379
func NewRedisConnection() (redis.Conn, error) {
	return redis.Dial("tcp", ":6379")
}

func (rs *RedisSet) Add(value string) error {
	_, err := rs.c.Do("SADD", rs.id, value)
	return err
}

func (rs *RedisSet) Has(value string) (bool, error) {
	return redis.Bool(rs.c.Do("SISMEMBER", rs.id, value))
}

func (rs *RedisSet) GetAll() ([]string, error) {
	result, err := redis.Values(rs.c.Do("SMEMBERS", rs.id))
	strs := make([]string, len(result))
	for i := 0; i < len(result); i++ {
		strs[i] = getString(result, i)
	}
	return strs, err

}

func (rs *RedisSet) Del(value string) error {
	_, err := rs.c.Do("SREM", rs.id, value)
	return err
}

func (rl *RedisList) Store(value string) error {
	_, err := rl.c.Do("RPUSH", rl.id, value)
	return err
}

func (rm *RedisKeyValue) Set(key, value string) error {
	_, err := rm.c.Do("SET", rm.id+":"+key, value)
	return err
}

func (rm *RedisKeyValue) Get(key string) (string, error) {
	result, err := redis.String(rm.c.Do("GET", rm.id+":"+key))
	if err != nil {
		return "", err
	}
	return result, nil
}

func (rh *RedisHashMap) Set(hashkey, key, value string) error {
	_, err := rh.c.Do("HSET", rh.id+":"+hashkey, key, value)
	return err
}

func (rh *RedisHashMap) Get(hashkey, key string) (string, error) {
	result, err := redis.String(rh.c.Do("HGET", rh.id+":"+hashkey, key))
	if err != nil {
		return "", err
	}
	return result, nil
}

func (rh *RedisHashMap) Del(hashkey, key string) error {
	_, err := rh.c.Do("HDEL", rh.id+":"+hashkey, key)
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
