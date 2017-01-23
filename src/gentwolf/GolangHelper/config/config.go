package config

import (
	"encoding/json"
	"io/ioutil"
)

func LoadConfig(filename string) (Config, error) {
	cfg := Config{}

	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return cfg, err
	}

	err = json.Unmarshal(b, &cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}

type WebConfig struct {
	Port    string `json:"port"`
	IsDebug bool   `json:"isDebug"`
}

type DbConfig struct {
	Type               string `json:"type"`
	Dsn                string `json:"dsn"`
	MaxOpenConnections int    `json:"maxOpenConnections"`
	MaxIdleConnections int    `json:"maxIdleConnections"`
}

type CacheConfig struct {
	Expiration int32  `json:"expiration"`
	Prefix     string `json:"prefix"`
	Host       string `json:"host"`
}

type Config struct {
	Web   WebConfig           `json:"web"`
	Db    map[string]DbConfig `json:"db"`
	Cache CacheConfig         `json:"cache"`
}
