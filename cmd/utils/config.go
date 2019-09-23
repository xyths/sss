package utils

import (
	"encoding/json"
	"log"
	"os"
)

type MongoConf struct {
	Host        string `json:"host"`
	Port        string `json:"port"`
	Database    string `json:"database"`
	MaxPoolSize uint64 `json:"maxPoolSize"`
	MinPoolSize uint64 `json:"minPoolSize"`
	AppName     string `json:"appName"`
}

func (mc MongoConf) URI() string {
	return "mongodb://" + mc.Host + mc.Port
}

type Mail struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
	BCC      string `json:"bcc"`
	Subject  string `json:"subject"`
}

type Template struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

type MailConfig struct {
	Mongo    MongoConf `json:"mongo,omitempty"`
	Mail     Mail      `json:"mail,omitempty"`
	Template Template  `json:"template,omitempty"`
}

func Load(file string) *MailConfig {
	var config MailConfig
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		log.Println(err)
	}
	err = json.NewDecoder(configFile).Decode(config)
	if err != nil {
		log.Println(err)
		return nil
	}
	return &config
}

type AppendConfig struct {
	Mongo MongoConf `json:"mongo"`
}

func LoadAppendConfig(file string) *AppendConfig {
	var config AppendConfig
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		log.Println(err)
	}
	err = json.NewDecoder(configFile).Decode(config)
	if err != nil {
		log.Println(err)
		return nil
	}
	return &config
}
