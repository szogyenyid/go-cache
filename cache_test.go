package cache

import (
	"sync"
	"testing"
	"time"
)

func TestKeyValueStore(t *testing.T) {
	t.Run("Put and Get", func(t *testing.T) {
		kvStore := NewKeyValueStore()
		kvStore.Put("key1", "value1", time.Second*1)
		kvStore.Put("key2", "value2", time.Second*1)

		value, found := kvStore.Get("key1")
		if !found || value != "value1" {
			t.Errorf("Failed to retrieve value for key1")
		}

		value, found = kvStore.Get("key2")
		if !found || value != "value2" {
			t.Errorf("Failed to retrieve value for key2")
		}
	})

	t.Run("Expiration", func(t *testing.T) {
		kvStore := NewKeyValueStore()
		kvStore.Put("key1", "value1", time.Millisecond*10)
		kvStore.Put("key2", "value2", time.Millisecond*50)

		time.Sleep(time.Millisecond * 10)

		value, found := kvStore.Get("key1")
		if found || value != nil {
			t.Errorf("Expired key1 was not deleted")
		}

		value, found = kvStore.Get("key2")
		if !found || value != "value2" {
			t.Errorf("Unexpected result for key2 after expiration")
		}
	})

	t.Run("Delete", func(t *testing.T) {
		kvStore := NewKeyValueStore()
		kvStore.Delete("key2")

		value, found := kvStore.Get("key2")
		if found || value != nil {
			t.Errorf("Deleted key2 was still found")
		}
	})

	t.Run("Concurrent Put and Get", func(t *testing.T) {
		kvStore := NewKeyValueStore()
		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()
			kvStore.Put("key3", "value3", time.Second*3)
		}()
		time.Sleep(time.Millisecond * 3)
		go func() {
			defer wg.Done()
			value, found := kvStore.Get("key3")
			if !found || value != "value3" {
				t.Errorf("Failed to retrieve value for key3")
			}
		}()

		wg.Wait()
	})

	t.Run("Concurrent Put and Delete", func(t *testing.T) {
		kvStore := NewKeyValueStore()
		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()
			kvStore.Put("key4", "value4", time.Second*2)
		}()

		go func() {
			defer wg.Done()
			time.Sleep(time.Millisecond * 10)
			kvStore.Delete("key4")
		}()

		wg.Wait()

		value, found := kvStore.Get("key4")
		if found || value != nil {
			t.Errorf("Deleted key4 was still found")
		}
	})

	t.Run("Expiration with concurrent access", func(t *testing.T) {
		kvStore := NewKeyValueStore()
		kvStore.Put("key5", "value5", time.Duration(0))

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()
			time.Sleep(time.Millisecond * 20)
			value, found := kvStore.Get("key5")
			if found || value != nil {
				t.Errorf("Expired key5 was not deleted")
			}
		}()

		go func() {
			defer wg.Done()
			time.Sleep(time.Millisecond * 10)
			kvStore.Delete("key5")
		}()

		wg.Wait()
	})

	t.Run("Get for non-existing key", func(t *testing.T) {
		kvStore := NewKeyValueStore()
		value, found := kvStore.Get("nonexistent")
		if found || value != nil {
			t.Errorf("Unexpected result for non-existing key")
		}
	})
}
