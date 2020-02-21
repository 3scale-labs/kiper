package threescale

import (
	"strconv"
	"time"

	"github.com/go-redis/redis/v7"
)

type redisLimitsStorage struct {
	internalStorage *redis.Client
}

func newRedisLimitsStorage(redisURL string) *redisLimitsStorage {
	opts := &redis.Options{
		Addr: redisURL,
	}
	return &redisLimitsStorage{internalStorage: redis.NewClient(opts)}
}

func (storage *redisLimitsStorage) get(key string) (int, bool, error) {
	val, err := storage.internalStorage.Get(key).Result()

	if err != nil {
		if err == redis.Nil {
			return 0, false, nil
		}
		return 0, false, err
	}

	valInt, err := strconv.Atoi(val)
	if err != nil {
		return 0, false, err
	}

	return valInt, true, nil
}

func (storage *redisLimitsStorage) create(key string, value int, duration time.Duration) (bool, error) {
	res, err := storage.internalStorage.SetNX(key, value, duration).Result()

	if err != nil {
		return false, err
	}

	return res, nil
}

func (storage *redisLimitsStorage) decrement(key string, value int) error {
	_, err := storage.internalStorage.DecrBy(key, int64(value)).Result()
	return err
}
