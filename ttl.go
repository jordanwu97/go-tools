package gotools

import (
	"fmt"
	"sync"
	"time"
)

// TTL time to live map. Don't use zero value.
type TTL struct {
	items    map[interface{}]*time.Timer
	itemsMtx sync.RWMutex
	expired  chan interface{}
}

// GetItems returns all non-expired items
func (ttl *TTL) GetItems() []interface{} {
	ttl.itemsMtx.RLock()
	defer ttl.itemsMtx.RUnlock()
	s := make([]interface{}, 0, len(ttl.items))
	for item := range ttl.items {
		s = append(s, item)
	}
	return s
}

// Expired channel delivers items that expires
func (ttl *TTL) Expired() <-chan interface{} {

	if ttl.items == nil {
		panic("ttl not instantiated")
	}
	if ttl.expired == nil {
		panic("ttl was set to not deliver expired items")
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

	ttl.itemsMtx.Lock()
	defer ttl.itemsMtx.Unlock()
	if t, exist := ttl.items[item]; exist {
		// reset timer for item if it exists
		t.Stop()
		t.Reset(expireIn)
	} else {
		// start a new timer for item. when timer expires, stop it, delete it, and deliver it on expired channel
		ttl.items[item] = time.AfterFunc(expireIn, func() {
			ttl.itemsMtx.Lock()
			ttl.items[item].Stop()
			delete(ttl.items, item)
			ttl.itemsMtx.Unlock()
			if ttl.expired != nil {
				ttl.expired <- item
			}
		})
	}

}

// CheckItem returns whether the item is in the TTL
func (ttl *TTL) CheckItem(item interface{}) bool {
	ttl.itemsMtx.RLock()
	defer ttl.itemsMtx.RUnlock()
	_, exist := ttl.items[item]
	return exist
}

// NewTTL instantiates the ttl. deliverExpired chooses whether to deliver expired items on ttl.Expired() channel. If set to false, calling ttl.Expired() will panic.
func NewTTL(deliverExpired bool) *TTL {
	ttl := TTL{
		items:   make(map[interface{}]*time.Timer),
		expired: nil,
	}
	if deliverExpired {
		ttl.expired = make(chan interface{})
	}
	return &ttl
}
