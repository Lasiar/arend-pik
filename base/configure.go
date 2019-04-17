package base

import (
	"encoding/json"
	"log"
	"os"
	"sync"
)

// Config basic config
type Config struct {
	ConnStr string `json:"connect_string"`
	Port    string `json:"port"`
}

var (
	_config     *Config
	_onceConfig sync.Once
)

// GetConfig get singleton config
func GetConfig() *Config {
	_onceConfig.Do(func() {
		_config = new(Config)
		_config.load()
	})
	return _config
}

func (c *Config) load() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	confFile, err := os.Open("/etc/pik-arenda/config.json")
	if err != nil {
		log.Fatal(err)
	}

	dc := json.NewDecoder(confFile)
	if err := dc.Decode(&c); err != nil {
		log.Panicf("Read Config file: %v", err)
	}

	if c.ConnStr == "" {
		log.Panicf("Can`t read connection string: %v", c.ConnStr)
	}
}
