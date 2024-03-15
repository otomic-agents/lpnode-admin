package redisbus

import (
	"admin-panel/redis_database"
	"log"
	"sync"

	"github.com/tidwall/gjson"
)

var redisBusOnce sync.Once

var redisBusIns *RedisBus

type RedisBus struct {
	redisDB *redis_database.RedisDb
}

func GetRedisBus() *RedisBus {
	redisBusOnce.Do(func() {
		redisBusIns = &RedisBus{}
		redisDB := redis_database.NewRedis("main")
		redisBusIns.redisDB = redisDB
	})
	return redisBusIns
}

func (rb *RedisBus) PublishEvent(channel string, val string) {
	go func() {
		event := gjson.Get(val, "type").String()
		payload := gjson.Get(val, "payload").Raw
		log.Println("publish event", channel, event, payload)
		rb.redisDB.Publish(channel, val)
	}()
}
