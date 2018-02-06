// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

import "time"

type Config struct {
	Period    time.Duration       `config:"period"`
	Host      string              `config:"host"`
	OIDs      []string            `config:"oids"`
	Community string              `config:"community"`
	Version   string              `config:"version"`
	Names     []map[string]string `config:"names"`
}

var DefaultConfig = Config{
	Period:    1 * time.Second,
	Host:      "127.0.0.1",
	Community: "public",
	Version:   "2c",
}
