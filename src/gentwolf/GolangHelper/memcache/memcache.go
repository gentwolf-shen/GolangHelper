package memcache

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"gentwolf/GolangHelper/config"
	MC "github.com/bradfitz/gomemcache/memcache"
	"time"
)

var (
	client            *MC.Client
	defaultExpiration = int32(60)
	keyPrefix         = "mem_"
)

func Init(cfg config.CacheConfig) error {
	return Connect(cfg.Expiration, cfg.Prefix, cfg.Host)
}

func Connect(expiration int32, prefix string, host ...string) error {
	selector := new(MC.ServerList)
	err := selector.SetServers(host...)
	if err != nil {
		return err
	}

	keyPrefix = prefix

	client = MC.NewFromSelector(selector)
	client.Timeout = 200 * time.Millisecond
	defaultExpiration = expiration

	return nil
}

func GetConn() *MC.Client {
	return client
}

func GetItem(name string) (*MC.Item, error) {
	return client.Get(keyPrefix + name)
}

func SetCAS(item *MC.Item) error {
	return client.CompareAndSwap(item)
}

func SetBytes(name string, bytes []byte, args ...int32) error {
	expire := getExpire(args...)
	return client.Set(&MC.Item{Key: keyPrefix + name, Value: bytes, Expiration: expire})
}

func GetBytes(name string) ([]byte, error) {
	item, err := client.Get(keyPrefix + name)
	if err != nil {
		return nil, err
	}

	return item.Value, nil
}

func Set(name string, value interface{}, args ...int32) error {
	bytes, err := encode(value)
	if err != nil {
		return err
	}

	expire := getExpire(args...)
	return client.Set(&MC.Item{Key: keyPrefix + name, Value: bytes, Expiration: expire})
}

func Get(name string, value interface{}) error {
	item, err := client.Get(keyPrefix + name)
	if err != nil {
		return err
	}

	return decode(item.Value, value)
}

func Delete(name string) error {
	return client.Delete(keyPrefix + name)
}

func encode(value interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(value)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func decode(b []byte, v interface{}) error {
	dec := gob.NewDecoder(bytes.NewBuffer(b))
	return dec.Decode(v)
}

func SetObject(name string, value interface{}, args ...int32) error {
	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}

	expire := getExpire(args...)
	return client.Set(&MC.Item{Key: keyPrefix + name, Value: bytes, Expiration: expire})
}

func GetObject(name string, value interface{}) error {
	item, err := client.Get(keyPrefix + name)
	if err != nil {
		return err
	}
	return json.Unmarshal(item.Value, value)
}

func getExpire(args ...int32) int32 {
	expire := defaultExpiration
	if len(args) > 0 {
		expire = args[0]
	}

	return expire
}

func Increment(key string, delta uint64) (uint64, error) {
	return client.Increment(keyPrefix+key, delta)
}

func Decrement(key string, delta uint64) (uint64, error) {
	return client.Decrement(keyPrefix+key, delta)
}
