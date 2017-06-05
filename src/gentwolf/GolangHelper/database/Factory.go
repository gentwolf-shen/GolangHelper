package database

import (
	"gentwolf/GolangHelper/config"
)

type factory struct {
	drivers map[string]*Base
}


var Factory *factory

func init() {
	Factory = &factory{}
}


func (this *factory) Init(configs map[string]config.DbConfig) error {
	this.drivers = make(map[string]*Base, len(configs))

	for k, cfg := range configs {
		tmp := Base{}
		err := tmp.OpenDb(cfg.Type, cfg.Dsn, cfg.MaxOpenConnections, cfg.MaxIdleConnections)
		if err != nil {
			return err
		} else {
			this.drivers[k] = &tmp
		}
	}
	return nil
}

func (this *factory) Driver(key string) *Base {
	return this.drivers[key]
}

func (this *factory) Close(key string) {
	this.drivers[key].Close()
}

func (this *factory) CloseAll() {
	for key, _ := range this.drivers {
		this.drivers[key].Close()
	}
}
