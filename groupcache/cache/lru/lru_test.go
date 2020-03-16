package lru

import (
	"reflect"
	"testing"
)

type String string

// Implement Value interface
func (d String) Len() int {
	return len(d)
}

func TestGet(t *testing.T) {
	lru := New(int64(0), nil)
	lru.Add("key1", String("1234"))
	if v, ok := lru.Get("key1"); !ok || string(v.(String)) != "1234" {
		t.Fatalf("Cache hit key1=1234 failed")
	}
	if _, ok := lru.Get("key2"); ok {
		t.Fatalf("Cache miss key2 failed")
	}

}

func TestRemoveOldest(t *testing.T) {
	k1, k2, k3, k4 := "k1", "k2", "k3", "k4"
	v1, v2, v3, v4 := "v1", "v2", "v3", "v4"

	capacity := len(k1 + k2 + k3 + v1 + v2 + v3)
	lru := New(int64(capacity), nil)
	lru.Add(k1, String(v1))
	lru.Add(k2, String(v2))
	lru.Add(k3, String(v3))
	lru.Add(k4, String(v4))

	if _, ok := lru.Get(k1); ok || lru.Len() != 3 {
		t.Fatalf("key1 is failed to remove")
	}

}

func TestRemove(t *testing.T) {
	removeKey := make([]string, 0)
	callback := func(key string, value Value) {
		removeKey = append(removeKey, key)
	}
	k1, k2, k3, k4 := "k1", "k2", "k3", "k4"
	v1, v2, v3, v4 := "v1", "v2", "v3", "v4"

	capacity := len(k1 + k2 + k3 + v1 + v2 + v3)
	lru := New(int64(capacity), callback)
	lru.Add(k1, String(v1))
	lru.Add(k2, String(v2))
	lru.Add(k3, String(v3))
	lru.Add(k4, String(v4))

	expect := []string{k1}
	if !reflect.DeepEqual(expect, removeKey) {
		t.Fatalf("call Remove function failed, expect %s\n but get %s\n", expect, removeKey)
	}

}
