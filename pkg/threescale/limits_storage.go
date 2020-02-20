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

// Returns true if the key has been created. False otherwise.
func (storage *LimitsStorage) create(key string, value int, duration time.Duration) bool {
	alreadyExistsErr := storage.internalStorage.Add(key, value, duration)
	return alreadyExistsErr == nil
}

func (storage *LimitsStorage) decrement(key string, value int) error {
	_, err := storage.internalStorage.DecrementInt(key, value)
	return err
}
