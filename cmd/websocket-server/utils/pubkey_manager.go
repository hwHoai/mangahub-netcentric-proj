package utils

import (
	"crypto/rsa"
	"sync"
)

var (
	publicKey *rsa.PublicKey
	mu        sync.RWMutex
)

func SetPublicKey(pub *rsa.PublicKey) {
	mu.Lock()
	defer mu.Unlock()
	publicKey = pub
}

func GetPublicKey() *rsa.PublicKey {
	mu.RLock()
	defer mu.RUnlock()
	return publicKey
}
