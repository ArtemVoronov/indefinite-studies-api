package cache

import (
	"hash/maphash"
	"sync"
)

const (
	BUCKET_NUMBER = 16
)

type Cache interface {
	Get(k string) (string, bool)
	Set(k string, v string)
}

type SimpleCache struct {
	store map[string]string
}

func (p *SimpleCache) Get(k string) (string, bool) {
	result, ok := p.store[k]
	return result, ok
}

func (p *SimpleCache) Set(k string, v string) {
	p.store[k] = v
}

type ConcurrentCache struct {
	stores [BUCKET_NUMBER]map[string]string
	rws    [BUCKET_NUMBER]sync.RWMutex
}

func getBucketNumber(k string) int {
	var h maphash.Hash
	h.WriteString(k)
	result := h.Sum64() % BUCKET_NUMBER
	return int(result)
}

func (p *ConcurrentCache) Get(k string) (string, bool) {
	bucketNumber := getBucketNumber(k)
	p.rws[bucketNumber].RLock()
	val, ok := p.stores[bucketNumber][k]
	p.rws[bucketNumber].RUnlock()
	return val, ok
}

func (p *ConcurrentCache) Set(k string, v string) {
	bucketNumber := getBucketNumber(k)
	p.rws[bucketNumber].Lock()
	p.stores[bucketNumber][k] = v
	p.rws[bucketNumber].Unlock()
}

func CreateSimpleCache() *SimpleCache {
	return &SimpleCache{store: make(map[string]string)}
}

func CreateConcurrentCache() *ConcurrentCache {
	result := &ConcurrentCache{}

	for i := 0; i < BUCKET_NUMBER; i++ {
		result.stores[i] = make(map[string]string)
	}

	return result
}
