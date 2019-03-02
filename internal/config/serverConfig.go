package config

import (
	"encoding/json"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"time"
)

// ServerProxyConfig contains parameters for the proxy server
type ServerProxyConfig struct {
	Enabled           bool          `yaml:"Enabled"`
	Scheme            string        `yaml:"Scheme"`
	Address           string        `yaml:"Address"`
	Path              string        `yaml:"Path"`
	User              string        `yaml:"User"`
	Password          string        `yaml:"Password"`
	ReconnectInterval time.Duration `yaml:"ReconnectInterval"`
}

// ConnectionConfig contains parameters for multicast -> websocket connections
type ConnectionConfig struct {
	SubscribePath    string            `yaml:"SubscribePath"`
	SendingInterval  time.Duration     `yaml:"SendingInterval"`
	MulticastAddress string            `yaml:"MulticastAddress"`
	ServerProxy      ServerProxyConfig `yaml:"ServerProxy"`
}

// RefereeConnection contains referee specific connection parameters
type RefereeConnection struct {
	ConnectionConfig `yaml:"Connection"`
}

// VisionConnection contains vision specific connection parameters
type VisionConnection struct {
	GeometrySendingInterval time.Duration `yaml:"GeometrySendingInterval"`
	ConnectionConfig        `yaml:"Connection"`
}

// ServerConfig is the root config containing all configs for the server
type ServerConfig struct {
	ListenAddress     string            `yaml:"ListenAddress"`
	RefereeConnection RefereeConnection `yaml:"RefereeConnection"`
	VisionConnection  VisionConnection  `yaml:"VisionConnection"`
}

func (s ServerConfig) String() string {
	str, err := json.Marshal(s)
	if err != nil {
		return err.Error()
	}
	return string(str)
}

// ReadServerConfig reads the server config from a yaml file
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
