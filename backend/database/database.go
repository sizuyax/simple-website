package database

import (
	"fmt"
	"github.com/go-redis/redis"
)

var client *redis.Client

func InitDB() {
	client = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})

	_, err := client.Ping().Result()
	if err != nil {
		fmt.Println("Ошибка подключения к Redis:", err)
		return
	}
	fmt.Println("Подключение к Redis успешно")
}

func GetClient() *redis.Client {
	return client
}
