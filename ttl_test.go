package gotools

import (
	"testing"
	"time"
)

func TestTTLFail(t *testing.T) {

	func() {
		ttl := &TTL{}
		defer func() {
			if r := recover(); r != nil {
				t.Log("recovered", r)
			}
		}()
		<-ttl.Expired()
	}()

	func() {
		ttl := &TTL{}
		defer func() {
			if r := recover(); r != nil {
				t.Log("recovered", r)
			}
		}()
		ttl.AddItem("1", 2*time.Second)
	}()

	func() {
		ttl := NewTTL()
		defer func() {
			if r := recover(); r != nil {
				t.Log("recovered", r)
			}
		}()
		ttl.AddItem("1", time.Nanosecond)
	}()

}

func TestTTL(t *testing.T) {

	ttl := NewTTL()
	ttl.AddItem("1", 2*time.Second)
	t1 := time.Now()
	ttl.AddItem("2", 3*time.Second)
	t2 := time.Now()

	time.Sleep(time.Second)
	ttl.AddItem("1", 5*time.Second)

	thresh := 10 * time.Millisecond

	if v := <-ttl.Expired(); v != "2" {
		t.Fatalf("Incorrect value. Want %v, got %v", "2", v)
	}
	if delay := time.Since(t2); delay-3*time.Second > thresh {
		t.Fatalf("Out of bound. Want %v, got %v", 3*time.Second, delay)
	}
	t.Log(time.Since(t2))

	if v := <-ttl.Expired(); v != "1" {
		t.Fatalf("Incorrect value. Want %v, got %v", "1", v)
	}
	if delay := time.Since(t1); delay-6*time.Second > thresh {
		t.Fatalf("Out of bound. Want %v, got %v", 6*time.Second, delay)
	}
	t.Log(time.Since(t2))
}
