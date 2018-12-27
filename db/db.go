package db

import (
	"estimate/models"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/websocket"
)

// https://github.com/pete911/examples-redigo

type WSUserMap = map[string]map[*websocket.Conn]bool

type Store struct {
	Users     WSUserMap
	Broadcast chan models.Estimation
	sync.Mutex
}

var (
	Pool       *redis.Pool
	wsStore    *Store
	pubSubConn *redis.PubSubConn
)

func Ping() error {

	conn := Pool.Get()
	defer conn.Close()

	_, err := redis.String(conn.Do("PING"))
	if err != nil {
		panic(err)
	}
	return nil
}

func newPool(server string) *redis.Pool {

	return &redis.Pool{

		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,

		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			return c, err
		},

		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

func Init() {
	redisHost := os.Getenv("REDIS_HOST")
	wsStore = &Store{
		Users:     make(WSUserMap),
		Broadcast: make(chan models.Estimation),
	}
	if redisHost == "" {
		redisHost = ":6379"
	}
	Pool = newPool(redisHost)
	Ping()
}

func Get(key string) ([]byte, error) {

	conn := Pool.Get()
	defer conn.Close()

	var data []byte
	data, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return data, fmt.Errorf("error getting key %s: %v", key, err)
	}
	return data, err
}

func GetString(key string) string {
	str, err := redis.String(Get(key))
	if err != nil {
		panic(err)
	}
	return str
}

func Set(key string, value []byte) error {
	conn := Pool.Get()
	defer conn.Close()

	_, err := conn.Do("SET", key, value)
	if err != nil {
		v := string(value)
		if len(v) > 15 {
			v = v[0:12] + "..."
		}
		return fmt.Errorf("error setting key %s to %s: %v", key, v, err)
	}
	return err
}

func Exists(key string) (bool, error) {

	conn := Pool.Get()
	defer conn.Close()

	ok, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return ok, fmt.Errorf("error checking if key %s exists: %v", key, err)
	}
	return ok, err
}

func Delete(key string) error {

	conn := Pool.Get()
	defer conn.Close()

	_, err := conn.Do("DEL", key)
	return err
}
