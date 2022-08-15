package config

import (
	"fmt"
	"github.com/go-redis/redis/v7"
	"github.com/spf13/viper"
	"log"
	"os"
)

var (
	configFile = "config/env.conf.yml"
	configType = "yml"
)

type (
	Configuration struct {
		Debug          bool     `mapstructure:"debug"`
		ContextTimeout int      `mapstructure:"contextTimeout"`
		Server         Server   `mapstructure:"server"`
		Database       Database `mapstructure:"database"`
		Grpc           Grpc     `mapstructure:"grpc"`
	}

	Server struct {
		Address       string `mapstructure:"address"`
		AccessSecret  string `mapstructure:"access-secret"`
		RefreshSecret string `mapstructure:"refresh-secret"`
	}

	Database struct {
		Driver   string `mapstructure:"driver"`
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
		Name     string `mapstructure:"name"`
		Sslmode  string `mapstructure:"sslmode"`
	}

	Grpc struct {
		Host    string `mapstructure:"host"`
		Port    int    `mapstructure:"port"`
		Network string `mapstructure:"tcp"`
	}
)

func NewConfig() Configuration {
	conf := &Configuration{}

	err := viper.Unmarshal(conf)

	if err != nil {
		fmt.Printf("unable decode into config struct, %v", err)
	}
	return *conf
}

func InitConfig() {
	viper.SetConfigType(configType)
	viper.SetConfigFile(configFile)

	err := viper.ReadInConfig()

	if err != nil {
		fmt.Println(err.Error())
	}
}

func InitRedis() *redis.Client {

	dsn := os.Getenv("REDIS_DSN")
	if len(dsn) == 0 {
		dsn = "localhost:6379"
	}

	log.Printf("Redis listening on PORT: %s", dsn)

	client := redis.NewClient(&redis.Options{
		Addr: dsn, //redis port
	})
	_, err := client.Ping().Result()
	if err != nil {
		panic(err)
	}
	return client
}
