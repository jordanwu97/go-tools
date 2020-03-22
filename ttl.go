package gotools

import (
	"fmt"
	"time"
)

// TTL time to live map. Don't use zero value.
type TTL struct {
	items   map[interface{}]*time.Timer
	expired chan interface{}
}

// Expired channel delivers items that expires
func (ttl *TTL) Expired() <-chan interface{} {
	if ttl.expired == nil {
		panic("ttl not instantiated")
	}
	return ttl.expired
}

// AddItem adds an item to the ttl, which expires in expireIn.
// If item is already in the TTL, will update the expiring time.
func (ttl *TTL) AddItem(item interface{}, expireIn time.Duration) {

	if ttl.items == nil {
		panic("ttl not instantiated")
	}

	if expireIn <= time.Nanosecond {
		panic(fmt.Errorf("expireIn must be larger than 1 nanosecond"))
	}

	if t, exist := ttl.items[item]; exist {
		t.Stop()
		t.Reset(expireIn)
	} else {
		ttl.items[item] = time.AfterFunc(expireIn, func() {
			ttl.items[item].Stop()
			delete(ttl.items, item)
			ttl.expired <- item
		})
	}
}

// NewTTL instantiates the ttl map
func NewTTL() *TTL {
	ttl := TTL{
		items:   make(map[interface{}]*time.Timer),
		expired: make(chan interface{}),
	}
	return &ttl
}
