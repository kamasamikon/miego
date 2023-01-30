package xredis

import (
	"github.com/go-redis/redis"

	"github.com/kamasamikon/miego/conf"
	"github.com/kamasamikon/miego/klog"
)

func Client() *redis.Client {
	Addr := conf.Str("172.17.0.1:6379", "s:/db/redis/addr")
	Pass := conf.Str("", "s:/db/redis/pass")
	DB := conf.Int(0, "i:/db/redis/db")
	redisdb := redis.NewClient(&redis.Options{
		Addr:     Addr,
		Password: Pass,
		DB:       int(DB),
	})

	_, err := redisdb.Ping().Result()
	if err != nil {
		klog.E("%s", err.Error())
		return nil
	}
	return redisdb
}
