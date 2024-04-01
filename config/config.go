package config

import "github.com/tkanos/gonfig"

type Configuration struct {
	DB_HOST        string
	DB_PORT        int
	DB_USER        string
	DB_PASSWORD    string
	DB_NAME        string
	SECRET_TOKEN_A string
}

func GetConfig() Configuration {
	conf := Configuration{}
	gonfig.GetConf("config/config.json", &conf)
	return conf
}
