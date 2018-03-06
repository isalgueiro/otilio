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
			for _, host := range bt.config.Hosts {
				gosnmp.Default.Target = host
				gosnmp.Default.Port = bt.config.Port
				gosnmp.Default.Community = bt.config.Community
				gosnmp.Default.Version = bt.version
				if bt.version == gosnmp.Version3 {
					gosnmp.Default.SecurityModel = gosnmp.UserSecurityModel
					gosnmp.Default.SecurityParameters = &gosnmp.UsmSecurityParameters{
						UserName:                 bt.config.User,
						AuthenticationPassphrase: bt.config.AuthPassword,
						PrivacyPassphrase:        bt.config.PrivPassword,
						AuthenticationProtocol:   gosnmp.SHA,
						PrivacyProtocol:          gosnmp.DES,
					}
				}
				err := gosnmp.Default.Connect()
				if err != nil {
					logp.Critical("Can't connect to %s: %v", host, err.Error())
					return fmt.Errorf("Can't connect to %s", host)
				}
				defer gosnmp.Default.Conn.Close()
				r, err := gosnmp.Default.Get(bt.oids)
				if err != nil {
					logp.Err("Can't get oids for %v: %v", host, err.Error())
				} else {
					event := common.MapStr{
						"@timestamp": common.Time(time.Now()),
						"type":       b.Name,
						"snmp.host":  host,
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
}

// Stop stops the beater
func (bt *Otilio) Stop() {
	bt.client.Close()
	close(bt.done)
}
