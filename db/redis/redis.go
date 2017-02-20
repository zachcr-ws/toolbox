package redis

import (
	"errors"
	"github.com/ZachBergh/toolbox/code"
	"github.com/garyburd/redigo/redis"
	"log"
	"time"
)

var RClient RedisClient

type RedisClient struct {
	redisPool *redis.Pool
}

func newRedisPool(addr string, speed int) RedisClient {
	r := &redis.Pool{
		MaxIdle:     speed,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", addr)
			if err != nil {
				log.Println(err)
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			if err != nil {
				log.Println(err)
			}
			return err
		},
	}
	return RedisClient{
		redisPool: r,
	}
}

func (r *RedisClient) GetString(key string) (string, error) {
	rc := r.redisPool.Get()
	defer rc.Close()

	rc.Send("GET", key)
	rc.Flush()
	return redis.String(rc.Receive())
}

func (r *RedisClient) HGetString(key string, tp string) (string, error) {
	rc := r.redisPool.Get()
	defer rc.Close()

	rc.Send("HGET", key, tp)
	rc.Flush()
	return redis.String(rc.Receive())
}

func (r *RedisClient) Set(key string, tp string) error {
	rc := r.redisPool.Get()
	defer rc.Close()

	rc.Send("MULTI")
	rc.Send("SET", key, tp)
	_, err := rc.Do("EXEC")
	return err
}

func (r *RedisClient) HSet(key string, tp string, val string) error {
	rc := r.redisPool.Get()
	defer rc.Close()

	rc.Send("MULTI")
	rc.Send("HSET", key, tp, val)
	_, err := rc.Do("EXEC")
	return err
}

func (r *RedisClient) HSetStruct(key string, field string, v interface{}, ttl int) error {
	rc := r.redisPool.Get()
	defer rc.Close()

	body, err := code.GobGeneralEncoder(v)
	if err != nil {
		return err
	}

	rc.Send("MULTI")
	rc.Send("HSET", key, field, body)
	if ttl > 0 {
		rc.Send("EXPIRE", key, ttl)
	}
	return rc.Send("EXEC")
}

func (r *RedisClient) HGetStruct(key string, tp string, structType interface{}, structVal interface{}) error {
	rc := r.redisPool.Get()
	defer rc.Close()

	exist := r.redisCheckKeyExists(key)
	if !exist {
		return errors.New("Key isn't exist")
	}

	body, err := rc.Do("HGET", key, tp)
	if err != nil {
		return err
	}

	if body == nil {
		return errors.New("Null")
	}
	return code.GobGeneralDecoder(body.([]byte), structType, structVal)
}

func (r *RedisClient) HSetDir(dir, key string, v interface{}, ttl int) error {
	rc := r.redisPool.Get()
	defer rc.Close()

	body, err := code.GobGeneralEncoder(v)
	if err != nil {
		return err
	}

	rc.Send("MULTI")
	rc.Send("SET", dir+":"+key, body)
	if ttl > 0 {
		rc.Send("EXPIRE", dir+":"+key, ttl)
	}
	return rc.Send("EXEC")
}

func (r *RedisClient) HGetDir(dir, key string, structType interface{}, structVal interface{}) error {
	rc := r.redisPool.Get()
	defer rc.Close()

	exist := r.redisCheckKeyExists(dir + ":" + key)
	if !exist {
		return errors.New("Key isn't exist")
	}

	body, err := rc.Do("GET", dir+":"+key)
	if err != nil {
		return err
	}

	if body == nil {
		return errors.New("Null")
	}
	return code.GobGeneralDecoder(body.([]byte), structType, structVal)
}

func (r *RedisClient) HDelDir(dir, key string) error {
	rc := r.redisPool.Get()
	defer rc.Close()

	exist := r.redisCheckKeyExists(dir + ":" + key)
	if !exist {
		return errors.New("Key isn't exist")
	}

	_, err := rc.Do("DEL", dir+":"+key)
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisClient) redisCheckKeyExists(key string) bool {
	conn := r.redisPool.Get()
	defer conn.Close()
	res, err := conn.Do("EXISTS", key)
	if err != nil {
		log.Println(err)
		return false
	}
	if res.(int64) == 1 {
		return true
	} else {
		return false
	}
}
