package redis

import (
	"gentwolf/GolangHelper/config"
	RDS "github.com/garyburd/redigo/redis"
	"strings"
	"time"
)

var pool *RDS.Pool

//初始化redis连接信息
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

//发送命令：无返回值
func Send(cmdName string, args ...interface{}) error {
	conn := pool.Get()
	err := conn.Send(cmdName, args...)
	conn.Close()

	return err
}

//发送命令：有返回值
func Do(cmdName string, args ...interface{}) (interface{}, error) {
	conn := pool.Get()
	replay, err := conn.Do(cmdName, args...)
	conn.Close()

	return replay, err
}

//设置内容
func Set(key string, v interface{}) error {
	return Send("SET", key, v)
}

//取回内容
func Get(key string) ([]byte, error) {
	return RDS.Bytes(Do("GET", key))
}

//取回内容(返回字符串)
func GetStr(key string) (string, error) {
	return RDS.String(Do("GET", key))
}

//消息队列：发送消息
func MqSet(key string, v interface{}) error {
	return Send("LPUSH", key, v)
}

//消息队列：取回消息
func MqGet(key string) (interface{}, error) {
	return Do("RPOP", key)
}

//发布/订阅模式：发布消息
func Publish(channel string, v interface{}) error {
	return Send("PUBLISH", channel, v)
}

//发布/订阅模式：订阅消息
func Subscribe(callback func(string, string, []byte), channel ...string) {
	go func() {
		conn := pool.Get()
		defer conn.Close()

		pubSub := RDS.PubSubConn{Conn: conn}
		defer pubSub.Close()

		for _, c := range channel {
			if strings.Contains(c, "*") {
				pubSub.PSubscribe(c)
			} else {
				pubSub.Subscribe(c)
			}
		}

		for {
			switch n := pubSub.Receive().(type) {
			case RDS.Message:
				callback("", n.Channel, n.Data)
			case RDS.PMessage:
				callback(n.Pattern, n.Channel, n.Data)
			}
		}
	}()
}
