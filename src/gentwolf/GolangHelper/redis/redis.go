package redis

import (
	"gentwolf/GolangHelper/config"
	RDS "github.com/garyburd/redigo/redis"
	"time"
)

var pool *RDS.Pool

func Connect(cfg config.RedisConfig) {
	pool = &RDS.Pool{}
	pool.MaxActive = cfg.MaxActive
	pool.MaxIdle = cfg.MaxIdle
	pool.Wait = cfg.Wait
	pool.IdleTimeout = time.Duration(cfg.IdleTimeout)

	pool.Dial = func() (RDS.Conn, error) {
		c, err := RDS.Dial("tcp", cfg.Address)
		if err != nil {
			return nil, err
		}
		return c, nil
	}

}

func Send(cmdName string, args ...interface{}) error {
	conn := pool.Get()
	err := conn.Send(cmdName, args...)
	conn.Close()

	return err
}

func Do(cmdName string, args ...interface{}) (interface{}, error) {
	conn := pool.Get()
	replay, err := conn.Do(cmdName, args...)
	conn.Close()

	return replay, err
}

func Set(key string, v interface{}) error {
	return Send("SET", key, v)
}

func Get(key string) ([]byte, error) {
	return RDS.Bytes(Do("GET", key))
}

func GetStr(key string) (string, error) {
	return RDS.String(Do("GET", key))
}

func MqSet(key string, v interface{}) error {
	return Send("LPUSH", key, v)
}

func MqGet(key string) (interface{}, error) {
	return Do("RPOP", key)
}
