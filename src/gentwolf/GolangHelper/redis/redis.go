package redis

import (
	RDS "github.com/garyburd/redigo/redis"
)

var conn RDS.Conn

func Dial(address string) error {
	var err error
	conn, err = RDS.Dial("tcp", address)
	return err
}

func Close() {
	conn.Close()
}

func Get(key string) ([]byte, error) {
	return RDS.Bytes(conn.Do("GET", key))
}

func GetStr(key string) (string, error) {
	return RDS.String(conn.Do("GET", key))
}

func Set(key string, v interface{}) error {
	_, err := conn.Do("SET", key, v)
	return err
}

func Send(cmdName string, args ...interface{}) error {
	return conn.Send(cmdName, args...)
}
