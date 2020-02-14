package threescale

import (
	"time"

	gocache "github.com/patrickmn/go-cache"
)

// For now, we use gocache to store the limits. In order to implement shared
// limits between instances, we could use a DB like Redis.

type LimitsStorage struct {
	internalStorage *gocache.Cache
}

func newLimitsStorage() LimitsStorage {
	goCache := gocache.New(gocache.NoExpiration, time.Minute)
	return LimitsStorage{internalStorage: goCache}
}

func (storage *LimitsStorage) get(key string) (int, bool) {
	val, exists := storage.internalStorage.Get(key)

	if !exists {
		return 0, false
	}

	return val.(int), exists
}

func (storage *LimitsStorage) setWithTTL(key string, value int, duration time.Duration) {
	storage.internalStorage.Set(key, value, duration)
}

func (storage *LimitsStorage) decrement(key string, value int) error {
	_, err := storage.internalStorage.DecrementInt(key, value)
	return err
}
