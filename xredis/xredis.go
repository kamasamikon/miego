package xredis

import (
	"github.com/go-redis/redis"

	"miego/conf"
	"miego/klog"
)

func Client(Addr string, Pass string, DB int) *redis.Client {
	if Addr == "" {
		Addr = conf.SGet("db/redis/addr", "172.17.0.1:6379")
	}
	if Pass == "" {
		Pass = conf.S("db/redis/pass")
	}
	if DB < 0 {
		DB = int(conf.I("db/redis/db", 0))
	}
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
