package engine

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sync"
)

const COLLISION_MAX = 3

type dataRecord struct {
	key   string
	value string
}

type dataTable struct {
	data  map[string]dataRecord
	mutex sync.RWMutex
}

func CreateEngine() *dataTable {
	return &dataTable{
		data:  make(map[string]dataRecord),
		mutex: sync.RWMutex{},
	}
}

func (t *dataTable) Set(key string, value string) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	hash, _, _, err := t.getHashAndValue(key)
	if err != nil {
		return err
	}

	t.data[hash] = dataRecord{key, value}

	return nil
}

func (t *dataTable) Get(key string) (string, bool, error) {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	_, value, ok, err := t.getHashAndValue(key)

	return value, ok, err
}

func (t *dataTable) Delete(key string) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	hash, _, ok, err := t.getHashAndValue(key)
	if err != nil {
		return err
	}

	if ok {
		delete(t.data, hash)

		return nil
	}

	return fmt.Errorf("key %s is not found", key)
}

func (t *dataTable) IsSet(key string) (bool, error) {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	_, _, ok, err := t.getHashAndValue(key)
	if err != nil {
		return false, err
	}

	return ok, nil
}

func (t *dataTable) Keys() *[]string {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	var keys = make([]string, len(t.data))

	idx := 0
	for _, value := range t.data {
		keys[idx] = value.key
		idx++
	}

	return &keys
}

func (t *dataTable) getHashAndValue(key string) (string, string, bool, error) {
	for i := 0; i < COLLISION_MAX; i++ {
		hash := t.createHash(key, i)
		record, ok := t.data[hash]

		if !ok {
			return hash, "", false, nil
		}
		if ok && record.key == key {
			return hash, record.value, true, nil
		}
	}

	return "", "", false, fmt.Errorf("key %s: too many collisions", key)
}

func (t *dataTable) createHash(key string, idx int) string {
	h := sha256.New()
	finalKey := key
	if idx > 0 {
		finalKey += fmt.Sprint(idx)
	}
	h.Write([]byte(finalKey))

	return hex.EncodeToString(h.Sum(nil))
}
