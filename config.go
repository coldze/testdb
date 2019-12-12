package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/go-redis/redis"
)

type redisCfg struct {
	Address string `json:"address"`
	DB      int    `json:"db"`
}

type mysqlCfg struct {
	Host string `json:"host"`
	User string `json:"user"`
	DB   string `json:"db"`
}

type bindCfg struct {
	Ip   string `json:"ip"`
	Port int    `json:"port"`
}

type appCfg struct {
	redisPassword     string   `json:"-"`
	mysqlPassword     string   `json:"-"`
	MysqlConnection   mysqlCfg `json:"mysql"`
	CacheTtlSeconds   int      `json:"cache_ttl_seconds"`
	AppTimeoutSeconds int      `json:"app_timeout_seconds"`
	Redis             redisCfg `json:"redis"`
	Bind              bindCfg  `json:"bind"`
}

func (a *appCfg) GetRedisOptions() *redis.Options {
	return &redis.Options{
		Addr:     a.Redis.Address,
		Password: a.redisPassword,
		DB:       a.Redis.DB,
	}
}

func (a *appCfg) GetMysqlURI() string {
	return fmt.Sprintf("%s:%s@%s/%s", a.MysqlConnection.User, a.mysqlPassword, a.MysqlConnection.Host, a.MysqlConnection.DB)
}

func (a *appCfg) GetBind() string {
	return fmt.Sprintf("%s:%v", a.Bind.Ip, a.Bind.Port)
}

func (a *appCfg) GetAppTimeout() time.Duration {
	return time.Duration(a.AppTimeoutSeconds) * time.Second
}

func (a *appCfg) GetCacheTtl() time.Duration {
	return time.Duration(a.CacheTtlSeconds) * time.Second
}

func getConfig(filename string, redisPassword string, mysqlPassword string) (*appCfg, error) {
	cfgData, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	cfg := &appCfg{
		redisPassword: redisPassword,
		mysqlPassword: mysqlPassword,
	}
	err = json.Unmarshal(cfgData, cfg)
	return cfg, err
}
