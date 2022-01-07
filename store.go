package main

import (
	"errors"
	"github.com/sirupsen/logrus"
	"sync"
)

var store = struct {
	sync.RWMutex
	m map[string]string
}{m: make(map[string]string)}

//var store = make(map[string]string)
var ErrorNoSuchKey = errors.New("no such key")

func Put(key string, value string) error {
	store.Lock()
	defer store.Unlock()
	store.m[key] = value
	logStoreState()
	return nil
}

func Get(key string) (string, error) {
	store.RLock()
	value, ok := store.m[key]
	store.RUnlock()
	if !ok {
		return "", ErrorNoSuchKey
	}
	logStoreState()
	return value, nil
}

func Delete(key string) error {
	store.Lock()
	defer store.Unlock()

	delete(store.m, key)
	logStoreState()
	return nil
}

func logStoreState() {
	logrus.Info("Current store = ", store.m)
}
