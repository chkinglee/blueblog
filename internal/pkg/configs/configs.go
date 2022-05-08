// Package configs
// @Author      : lilinzhen
// @Time        : 2022/5/7 15:52:45
// @Description :
package configs

import (
	"blueblog/pkg/env"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"time"

	"github.com/spf13/viper"
)

var config = new(Config)

type Config struct {
	Server struct {
		HttpPort int `yaml:"httpPort"`
		GrpcPort int `yaml:"grpcPort"`
	} `yaml:"server"`

	Blues struct {
		Blueblog struct {
			Interface struct {
				Addr string `yaml:"addr"`
			} `yaml:"interface"`
			Service struct {
				Addr string `yaml:"addr"`
			} `yaml:"service"`
			Job struct {
				Addr string `yaml:"addr"`
			} `yaml:"job"`
			Task struct {
				Addr string `yaml:"addr"`
			} `yaml:"task"`
			Admin struct {
				Addr string `yaml:"addr"`
			} `yaml:"admin"`
		} `yaml:"blueblog"`
	} `yaml:"blues"`

	Logger struct {
		File  string `yaml:"file"`
		Level string `yaml:"level"`
	} `yaml:"logger"`

	MySQL struct {
		Read struct {
			Addr string `yaml:"addr"`
			User string `yaml:"user"`
			Pass string `yaml:"pass"`
			Name string `yaml:"name"`
		} `yaml:"read"`
		Write struct {
			Addr string `yaml:"addr"`
			User string `yaml:"user"`
			Pass string `yaml:"pass"`
			Name string `yaml:"name"`
		} `yaml:"write"`
		Base struct {
			MaxOpenConn     int           `yaml:"maxOpenConn"`
			MaxIdleConn     int           `yaml:"maxIdleConn"`
			ConnMaxLifeTime time.Duration `yaml:"connMaxLifeTime"`
		} `yaml:"base"`
	} `yaml:"mysql"`

	Redis struct {
		Addr         string `yaml:"addr"`
		Pass         string `yaml:"pass"`
		Db           int    `yaml:"db"`
		MaxRetries   int    `yaml:"maxRetries"`
		PoolSize     int    `yaml:"poolSize"`
		MinIdleConns int    `yaml:"minIdleConns"`
	} `yaml:"redis"`

	RabbitMQ struct {
		Url string `yaml:"url"`
	} `yaml:"rabbitmq"`

	URLToken struct {
		Secret         string        `yaml:"secret"`
		ExpireDuration time.Duration `yaml:"expireDuration"`
	} `yaml:"urlToken"`

	HashIds struct {
		Secret string `yaml:"secret"`
		Length int    `yaml:"length"`
	} `yaml:"hashids"`

	Language struct {
		Local string `yaml:"local"`
	} `yaml:"language"`
}

func Init() {
	viper.SetConfigName(fmt.Sprintf("%s-%s", env.Active().App(), env.Active().Value()))
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	if err := viper.Unmarshal(config); err != nil {
		panic(err)
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		if err := viper.Unmarshal(config); err != nil {
			panic(err)
		}
	})
}

func Get() Config {
	return *config
}
