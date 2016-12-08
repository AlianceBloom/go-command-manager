package command_manager

import (
	_ "gopkg.in/redis.v5"
	"gopkg.in/redis.v5"
)

var redisClient *redis.Client

func ConnectRedis() {
	redisClient = redis.NewClient(&redis.Options{
		Addr: RedisAddr,
		Password: RedisPassword,
		DB: 0,
	})

	_, err := redisClient.Ping().Result()

	PanicError(err)
}

func existInRedis(key string) bool {
	_, err := redisClient.Get(DefaultRedisKey + ":" + key).Result()
	if err == redis.Nil {
		return false

	} else if err != nil {
		PanicError(err)
	} else {
		return true
	}
	return false
}

func getFromRedisByKey(key string) string {
	result, err := redisClient.Get(DefaultRedisKey + ":" + key).Result()
	if err == redis.Nil {
		return ""
	} else if err != nil {
		PanicError(err)
	}

	return result
}

func deleteFromRedisByKey(key string) bool {
	err := redisClient.Del(DefaultRedisKey + ":" + key).Err()
	if err == redis.Nil {
		return false
	} else if err != nil {
		PanicError(err)
	}
	return true
}

func storeInRedis(key string, value interface{}) {
	err := redisClient.Set(DefaultRedisKey + ":" + key, value, 0).Err()
	PanicError(err)
}
