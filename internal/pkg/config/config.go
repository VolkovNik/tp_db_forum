package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type DBConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	DBName   string `json:"db,omitempty"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type Service struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

type Config struct {
	Postgres *DBConfig `json:"postgres"`
	Main     *Service  `json:"main"`
}

var config = &Config{}

func Get() *Config {
	return config
}

func init() {
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	err = json.Unmarshal(data, config)
	if err != nil {
		fmt.Println(err)
		return
	}
}
