package utils

import (
	"crypto/rsa"
	"sync"
)

var (
	PublicKey *rsa.PublicKey
	mu        sync.RWMutex
)

func SetPublicKey(pub *rsa.PublicKey) {
	mu.Lock()
	defer mu.Unlock()
	PublicKey = pub
}

func GetPublicKey() *rsa.PublicKey {
	mu.RLock()
	defer mu.RUnlock()
	return PublicKey
}
