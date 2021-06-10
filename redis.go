package datadog

import (
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/kiyor/env"
)

var Redis *RedisPool

type RedisPool struct {
	Pool redis.Pool
}

var redisHost string

func inDocker() bool {
	b, err := ioutil.ReadFile("/proc/1/cgroup")
	if err == nil && strings.HasSuffix(strings.Split(string(b), "\n")[0], "/") {
		return false
	}
	return true
}

func init() {
	env.StringVar(&redisHost, "REDIS_HOST", "127.0.0.1", redisHost)
	host := "127.0.0.1"
	if inDocker() {
		if redisHost == "127.0.0.1" {
			host = "redis"
		} else {
			host = redisHost
		}
	}
	Redis = &RedisPool{
		redis.Pool{
			MaxIdle:     6,
			IdleTimeout: 240 * time.Second,
			Dial: func() (redis.Conn, error) {
				c, err := redis.Dial("tcp", host+":6379")
				if err != nil {
					return nil, err
				}
				return c, err
			},
			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				_, err := c.Do("PING")
				return err
			},
		},
	}
}

func (r *RedisPool) Reset() {
	conn := r.Pool.Get()
	defer conn.Close()
	conn.Do("FLUSHALL")
}

func (r *RedisPool) Get(key string) ([]byte, bool) {
	conn := r.Pool.Get()
	defer conn.Close()
	res, err := conn.Do("GET", key)
	if err != nil {
		log.Println("redis", err)
		return nil, false
	}
	if res != nil {
		b := res.([]byte)
		return b, true
	} else {
		return nil, false
	}
}
func (r *RedisPool) GetWithTTL(key string) ([]byte, int64, bool) {
	conn := r.Pool.Get()
	defer conn.Close()
	res, err := conn.Do("GET", key)
	if err != nil {
		log.Println("redis", err)
		return nil, 0, false
	}
	if res != nil {
		b := res.([]byte)
		ttl, _ := conn.Do("TTL", key)
		return b, ttl.(int64), true
	} else {
		return nil, 0, false
	}
}

func (r *RedisPool) Set(key string, value []byte) error {
	conn := r.Pool.Get()
	defer conn.Close()
	_, err := conn.Do("SET", key, value)
	return err
}

func (r *RedisPool) SetWithTTL(key string, value []byte, second int) error {
	conn := r.Pool.Get()
	defer conn.Close()
	_, err := conn.Do("SET", key, value)
	if err != nil {
		return err
	}
	_, err = conn.Do("EXPIRE", key, second)
	return err
}
