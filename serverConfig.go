package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"time"
)

type ServerProxyConfig struct {
	Enabled           bool          `yaml:"Enabled"`
	Scheme            string        `yaml:"Scheme"`
	Address           string        `yaml:"Address"`
	Path              string        `yaml:"Path"`
	User              string        `yaml:"User"`
	Password          string        `yaml:"Password"`
	ReconnectInterval time.Duration `yaml:"ReconnectInterval"`
}

type ConnectionConfig struct {
	SubscribePath    string        `yaml:"SubscribePath"`
	SendingInterval  time.Duration `yaml:"SendingInterval"`
	MulticastAddress string        `yaml:"MulticastAddress"`
}

type ServerConfig struct {
	ServerProxy       ServerProxyConfig `yaml:"ServerProxy"`
	ListenAddress     string            `yaml:"ListenAddress"`
	RefereeConnection ConnectionConfig  `yaml:"RefereeConnection"`
	VisionConnection  ConnectionConfig  `yaml:"VisionConnection"`
}

func ReadServerConfig(fileName string) ServerConfig {
	config := ServerConfig{}
	f, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	d, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatalln(err)
	}
	err = yaml.Unmarshal(d, &config)
	if err != nil {
		log.Fatalln(err)
	}
	return config
}
