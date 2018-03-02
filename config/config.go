// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

import "time"

// Config stores Otilio configuration loaded from .yaml file
type Config struct {
	Period       time.Duration       `config:"period"`
	Host         string              `config:"host"`
	Port         uint16              `config:"port"`
	Community    string              `config:"community"`
	User         string              `config:"user"`
	AuthPassword string              `config:"authpass"`
	PrivPassword string              `config:"privpass"`
	Version      string              `config:"version"`
	OIDs         []map[string]string `config:"oids"`
}

// DefaultConfig default configuration
var DefaultConfig = Config{
	Period:    1 * time.Second,
	Host:      "127.0.0.1",
	Port:      161,
	Community: "public",
	Version:   "2c",
}
