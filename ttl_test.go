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

	if !ttl.CheckItem("1") {
		t.Fatalf("Item \"1\" should be in ttl")
	}

	if !ttl.CheckItem("2") {
		t.Fatalf("Item \"1\" should be in ttl")
	}

	time.Sleep(time.Second)
	ttl.AddItem("1", 5*time.Second)

	thresh := 10 * time.Millisecond

	if v := <-ttl.Expired(); v != "2" {
		t.Fatalf("Incorrect value. Want %v, got %v", "\"2\"", v)
	}
	if delay := time.Since(t2); delay-3*time.Second > thresh {
		t.Fatalf("Out of bound. Want %v, got %v", 3*time.Second, delay)
	}
	if ttl.CheckItem("2") {
		t.Fatalf("Item \"2\" should not be in ttl")
	}

	if v := <-ttl.Expired(); v != "1" {
		t.Fatalf("Incorrect value. Want %v, got %v", "\"1\"", v)
	}
	if delay := time.Since(t1); delay-6*time.Second > thresh {
		t.Fatalf("Out of bound. Want %v, got %v", 6*time.Second, delay)
	}
	if ttl.CheckItem("1") {
		t.Fatalf("Item \"1\" should not be in ttl")
	}
}
