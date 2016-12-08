package command_manager

import "flag"

var (
	RedisAddr string = *flag.String("redis-addr", "localhost:6379", "a string")
	RedisPassword string = *flag.String("redis-pass", "", "a string")
	DefaultRedisKey = *flag.String("default-redis-namespace", "command_manager", "string")
)


func PanicError(err error) {
	if err != nil {
		panic(err)
	}
}

func Run() {
	ConnectRedis()
	HandleHTTMRequests()
}