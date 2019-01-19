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

/*
{
	[sessionID1]: {
		[wsConn1]: true,
		[wsConn2]: true,
	},
	[sessionID2] {
		...
	}
*/
// globals

// WSUserMap holds all the users. TODO: add 'isAdmin' information
type WSUserMap = map[string]map[*websocket.Conn]bool

var (
	ClientSessions WSUserMap
)

type Store struct {
	Users     WSUserMap
	Broadcast chan models.UserMessageEstimation
	sync.Mutex
}

var (
	Pool    *redis.Pool
	WsStore *Store
	PSC     *redis.PubSubConn
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
	WsStore = &Store{
		Users:     make(WSUserMap),
		Broadcast: make(chan models.UserMessageEstimation),
	}
	if redisHost == "" {
		redisHost = ":6379"
	}
	Pool = newPool(redisHost)
	conn := Pool.Get()
	PSC = &redis.PubSubConn{Conn: conn}
	// whenever a key changes, we want to notify users.
	conn.Do("CONFIG", "SET", "notify-keyspace-events", "KEA")
	fmt.Println("Set the notify-keyspace-events to KEA")
	PSC.PSubscribe("__key*__:*")
	ClientSessions = WsStore.Users
	Ping()
}

/*

func (s *Store) findAndDeliver(userID string, content string) {
	m := Message{
		Content: content,
	}
	for _, u := range s.Users {
		if u.ID == userID {
			if err := u.conn.WriteJSON(m); err != nil {
				log.Printf("error on message delivery e: %s\n", err)
			} else {
				log.Printf("user %s found, message sent\n", userID)
			}
			return
		}
	}
	log.Printf("user %s not found at our store\n", userID)
}
*/

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
