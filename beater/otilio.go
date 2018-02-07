package beater

import (
	"fmt"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/publisher"
	"github.com/soniah/gosnmp"

	"github.com/isalgueiro/otilio/config"
)

// Otilio data type
type Otilio struct {
	done      chan struct{}
	config    config.Config
	client    publisher.Client
	version   gosnmp.SnmpVersion
	oidToName map[string]string
	oids      []string
}

// New creates beater
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	config := config.DefaultConfig
	if err := cfg.Unpack(&config); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	version := gosnmp.Version2c
	switch config.Version {
	case "1":
		version = gosnmp.Version1
	case "2c":
		version = gosnmp.Version2c
	case "3":
		version = gosnmp.Version3
	default:
		logp.Err("Wrong SNMP version %s, defaulting to 2c", config.Version)
	}

	m := make(map[string]string)
	var o []string
	for _, v := range config.OIDs {
		logp.Debug("otilio", "OID %s translates to %s in event", v["oid"], v["name"])
		m[v["oid"]] = v["name"]
		o = append(o, v["oid"])
	}

	bt := &Otilio{
		done:      make(chan struct{}),
		config:    config,
		version:   version,
		oidToName: m,
		oids:      o,
	}
	return bt, nil
}

// Run runs the beater
func (bt *Otilio) Run(b *beat.Beat) error {
	logp.Info("otilio is running! Hit CTRL-C to stop it.")

	bt.client = b.Publisher.Connect()
	ticker := time.NewTicker(bt.config.Period)
	for {
		select {
		case <-bt.done:
			return nil
		case <-ticker.C:
			// TODO: connect outside the loop with a timeout < bt.config.Period
			gosnmp.Default.Target = bt.config.Host
			err := gosnmp.Default.Connect()
			if err != nil {
				logp.Critical("Can't connect to %s: %v", bt.config.Host, err.Error())
				return fmt.Errorf("Can't connect to %s", bt.config.Host)
			}
			defer gosnmp.Default.Conn.Close()
			gosnmp.Default.Community = bt.config.Community
			gosnmp.Default.Version = bt.version
			r, err := gosnmp.Default.Get(bt.oids)
			if err != nil {
				logp.Err("Can't get oids %v: %v", bt.config.OIDs, err.Error())
			} else {
				event := common.MapStr{
					"@timestamp": common.Time(time.Now()),
					"type":       b.Name,
				}
				for _, v := range r.Variables {
					var value interface{}
					k := bt.oidToName[v.Name]
					if k == "" {
						k = v.Name
					}
					switch v.Type {
					case gosnmp.OctetString:
						value = string(v.Value.([]byte))
					default:
						value = gosnmp.ToBigInt(v.Value)
					}
					logp.Debug("otilio", "%s = %s", k, value)
					event.Put(k, value)
				}
				bt.client.PublishEvent(event)
				logp.Info("Event sent")
			}
		}
	}
}

// Stop stops the beater
func (bt *Otilio) Stop() {
	bt.client.Close()
	close(bt.done)
}
