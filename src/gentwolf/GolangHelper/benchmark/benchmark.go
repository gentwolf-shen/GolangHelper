package benchmark

import (
	"time"
)

var (
	items = make(map[string]time.Time)
)

func Start(name string) {
	items[name] = time.Now()
}

func stopTime(name string) time.Duration {
	item, bl := items[name]
	if !bl {
		return time.Now().Sub(time.Now())
	}
	delete(items, name)

	return time.Now().Sub(item)
}

func Stop(name string) float64 {
	return stopTime(name).Seconds()
}

func StopMillisecond(name string) float64 {
	return Stop(name) * 1000
}

func StopNano(name string) int64 {
	return stopTime(name).Nanoseconds()
}
