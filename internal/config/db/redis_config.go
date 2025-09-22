package db

import (
	"fmt"
	"k8s-platform-go/internal/config"

	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

func InitRedis() *redis.Client {
	redisProperties := config.GetGlobalConfig().Redis
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisProperties.Host, redisProperties.Port),
		Username: "",
		Password: redisProperties.Password, // no password set
		DB:       redisProperties.DB,       // use default DB
	})
	redisClient = client
	return redisClient
}

func CloseRedis() {
	err := redisClient.Close()
	if err != nil {
		fmt.Printf("redis close faild: %v", err)
		return
	}
}
