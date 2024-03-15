package database_config

import (
	"fmt"
	"log"
	"os"
)

// RedisDbConnectInfoItem Redis
type RedisDbConnectInfoItem struct {
	DbIndex  int
	RedisUrl string
	RedisPwd string
}
type RedisDataDataBaseConfig map[string]RedisDbConnectInfoItem

var RedisDataDataBaseConfigIns = make(map[string]RedisDbConnectInfoItem)

func Init() {
	InitMongoConfig()
	InitRedisConfig()
}

// InitRedisConfig init redis config required
func InitRedisConfig() {
	prodRedisHost := os.Getenv("REDIS_HOST")
	if prodRedisHost != "" {
		log.Println("use redis config from env vars")
		prodRedisPort := os.Getenv("REDIS_PORT")
		prodRedisPass := os.Getenv("REDIS_PASS")
		RedisDataDataBaseConfigIns["main"] = RedisDbConnectInfoItem{
			RedisUrl: fmt.Sprintf("%s:%s", prodRedisHost, prodRedisPort),
			RedisPwd: prodRedisPass,
			DbIndex:  0,
		}
		RedisDataDataBaseConfigIns["statusDb"] = RedisDbConnectInfoItem{
			RedisUrl: fmt.Sprintf("%s:%s", prodRedisHost, prodRedisPort),
			RedisPwd: prodRedisPass,
			DbIndex:  0,
		}
		return
	}
	redisPass := os.Getenv("REDIS_PASS")
	RedisDataDataBaseConfigIns["main"] = RedisDbConnectInfoItem{
		RedisUrl: "127.0.0.1:6379",
		RedisPwd: redisPass,
		DbIndex:  0,
	}
	RedisDataDataBaseConfigIns["statusDb"] = RedisDbConnectInfoItem{
		RedisUrl: "127.0.0.1:6379",
		RedisPwd: redisPass,
		DbIndex:  0,
	}
}
