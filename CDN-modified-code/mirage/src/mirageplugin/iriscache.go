package main

import "cache/eviction/iris"
import "cache/eviction/modifiedlru"

type CacheStorage struct {
	spectrumCache    *modifiedlru.CacheStorage
	mlruCache        *modifiedlru.CacheStorage
	serverSpectrum   uint64
	contentSpectrums []uint64
	spectrumCapacity int
	mlruCapacity     int
	hitCount         int
	missCount        int
}

func NewIrisCache(capacity int, spectrumRatio float64) iris.Accessor {
	irisCache := new(CacheStorage)
	return irisCache
}

func (cache *CacheStorage) CacheList() []interface{} {
	cacheList := make([]interface{}, 0)
	return cacheList
}

func (cache *CacheStorage) Inspect() {
}

func (cache *CacheStorage) FillUp() {
}

func (cache *CacheStorage) SetContentSpectrums(contentSpectrums []uint64) {
}

func (cache *CacheStorage) SetServerSpectrum(spectrum uint64) {
}

func (cache *CacheStorage) Len() int {
	return 0
}

func (cache *CacheStorage) Exist(key interface{}) bool {
	return false
}

func (cache *CacheStorage) Insert(key, value interface{}) interface{} {
	return nil
}

func (cache *CacheStorage) Fetch(key interface{}) interface{} {
	return nil
}

func (cache *CacheStorage) SpectrumCapacity() int {
	return 0
}

func (cache *CacheStorage) Capacity() int {
	return 0
}

func (cache *CacheStorage) HitCount() int {
	return 0
}

func (cache *CacheStorage) MissCount() int {
	return 0
}

func (cache *CacheStorage) ResetCount() {
}

func (cache *CacheStorage) Clear() {
}
