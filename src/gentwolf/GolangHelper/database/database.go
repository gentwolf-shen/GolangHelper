package database

import (
	"gentwolf/GolangHelper/config"
)

var (
	drivers map[string]*Base
)

func Init(configs map[string]config.DbConfig) error {
	drivers = make(map[string]*Base, len(configs))

	for k, cfg := range configs {
		tmp := Base{}
		err := tmp.OpenDb(cfg.Type, cfg.Dsn, cfg.MaxOpenConnections, cfg.MaxIdleConnections)
		if err != nil {
			return err
		} else {
			drivers[k] = &tmp
		}
	}
	return nil
}

func Driver(key string) *Base {
	return drivers[key]
}

func Close(key string) {
	drivers[key].Close()
}

func CloseAll() {
	for key, _ := range drivers {
		drivers[key].Close()
	}
}
