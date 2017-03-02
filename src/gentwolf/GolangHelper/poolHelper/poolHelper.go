package poolHelper

import (
	"sync"
	"time"
)

type Pool struct {
	New   func() interface{}
	Ping  func(interface{}) error
	Close func(interface{})

	store  chan interface{}
	min    int
	max    int
	locker sync.Mutex
}

func New(min, max int, pingTime int, newFunc func() interface{}) *Pool {
	p := &Pool{}
	p.max = max
	p.store = make(chan interface{}, max)
	p.New = newFunc

	go func() {
		for i := 1; i <= min; i++ {
			if v := p.create(); v != nil {
				p.store <- v
			}
		}
		go p.ping(pingTime)
	}()

	return p
}

func (this *Pool) create() interface{} {
	return this.New()
}

func (this *Pool) Get() interface{} {
	for {
		select {
		case v := <-this.store:
			return v
		default:
			return this.create()
		}
	}
}

func (this *Pool) Put(v interface{}) {
	select {
	case this.store <- v:
		return
	default:
		this.Close(v)
	}
}

func (this *Pool) Len() int {
	return len(this.store)
}

func (this *Pool) Destory() {
	this.locker.Lock()
	defer this.locker.Unlock()

	if this.store == nil {
		return
	}

	close(this.store)
	for v := range this.store {
		if v != nil {
			this.Close(v)
		}
	}

	this.store = nil
}

func (this *Pool) ping(pingTime int) {
	for _ = range time.Tick(time.Duration(int64(pingTime)) * time.Second) {
		length := this.Len()
		for i := 0; i < length; i++ {
			if v := this.Get(); v != nil {
				if err := this.Ping(v); err == nil {
					this.Put(v)
				}
			}
		}
	}
}
