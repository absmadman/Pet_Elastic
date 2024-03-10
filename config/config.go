package config

import (
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v8"
	"log"
	"os"
)

type Config struct {
	Username  string   `json:"username"`
	Password  string   `json:"password"`
	Addresses []string `json:"addresses"`
}

func NewCfg() *elasticsearch.Config {
	var cfg Config

	file, err := os.ReadFile("config/elasticsearch_user_cfg.json")

	if err != nil {
		log.Println(err)
	}

	err = json.Unmarshal(file, &cfg)

	if err != nil {
		log.Println(err)
	}

	return &elasticsearch.Config{
		Addresses: cfg.Addresses,
		Username:  cfg.Username,
		Password:  cfg.Password,
	}
}
