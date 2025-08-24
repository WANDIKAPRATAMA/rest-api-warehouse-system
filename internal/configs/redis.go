package configs

import (
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func NewRedis(viper *viper.Viper) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: viper.GetString("redis.host") + ":" + viper.GetString("redis.port"),
		// TODO: Filled with your redis password
		// Password: viper.GetString("redis.password"),
		DB: viper.GetInt("redis.db"),
	})
}
