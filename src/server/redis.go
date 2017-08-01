package main

import (
	"github.com/garyburd/redigo/redis"
	"time"
	"os"
	"os/signal"
	"syscall"
	"fmt"
)

var (
	Pool *redis.Pool
)

func init() {
	redisHost := ":6379"
	Pool = newPool(redisHost)
	close()
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

func close() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	signal.Notify(c, syscall.SIGKILL)
	go func() {
		<-c
		Pool.Close()
		os.Exit(0)
	}()
}

func Get(key string) ([]byte, error) {

	conn := Pool.Get()
	defer conn.Close()

	var data []byte
	data, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return data, fmt.Errorf("error the key %s is nil", key)
	}
	return data, err
}

func Set(key string, value string) (error) {
	conn := Pool.Get()
	defer conn.Close()

	_, err := conn.Do("SET", key, value)
	if err != nil {
		return fmt.Errorf("the err is: %v", err)
	}
	return nil
}

func Del(key string) (int, error) {
	conn := Pool.Get()
	defer conn.Close()

	_, err := conn.Do("DEL", key)
	if err != nil {
		return 0, fmt.Errorf("the err is: %v", err)
	}
	return 1, nil
}

//func main() {
//	err := Set("abc", "123")
//	if err != nil {
//		fmt.Println("set failed", err)
//	}
//	test, _ := Get("abc")
//	fmt.Println(string(test))
//}
