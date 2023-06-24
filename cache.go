package cache

import (
	"fmt"
	"sync"
	"time"
)

type keyValue struct {
	value      interface{}
	expiration *time.Time
}

type KeyValueStore struct {
	store map[string]keyValue
	mutex sync.RWMutex
}

func NewKeyValueStore() *KeyValueStore {
	return &KeyValueStore{
		store: make(map[string]keyValue),
		mutex: sync.RWMutex{},
	}
}

func (kv *KeyValueStore) Put(key string, value interface{}, expiration time.Duration) error {
	kv.mutex.Lock()
	defer kv.mutex.Unlock()

	var expirationTime *time.Time
	if expiration > 0 {
		expTime := time.Now().Add(expiration)
		expirationTime = &expTime
	}

	kv.store[key] = keyValue{
		value:      value,
		expiration: expirationTime,
	}

	return nil
}

func (kv *KeyValueStore) Get(key string) (interface{}, bool) {
	if item, found := kv.store[key]; found {
		if item.expiration != nil && time.Now().After(*item.expiration) {
			kv.mutex.Lock()
			defer kv.mutex.Unlock()

			// Double-check expiration after acquiring the write lock
			if item.expiration != nil && time.Now().After(*item.expiration) {
				delete(kv.store, key)
				return nil, false
			}
		}

		return item.value, true
	}

	return nil, false
}

func (kv *KeyValueStore) Delete(key string) error {
	kv.mutex.Lock()
	defer kv.mutex.Unlock()

	if _, found := kv.store[key]; !found {
		return fmt.Errorf("key not found: %s", key)
	}

	delete(kv.store, key)
	return nil
}
