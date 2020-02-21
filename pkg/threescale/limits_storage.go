package threescale

import (
	"time"
)

// Note: this interface is not the best for some storage backends. For example,
// in Redis, some of the operations could be performed in a single MULTI command
// to save network round-trips or several limits could be updated in a single
// pipeline. Something to think about in the future if a more efficient solution
// is needed.

type limitsStorage interface {
	// The first element returned is the value of the key. The second is a
	// boolean that indicates whether the key exists or not.
	get(key string) (int, bool, error)

	// Creates a key with the given values and ttl, only when it does not exist.
	// Returns true when the key was created because it did not exist before.
	// Returns false if the key was not created because it was already set.
	create(key string, value int, duration time.Duration) (bool, error)

	// Decreases the value of a key by the given number. Returns an error if the
	// key does not exist.
	decrement(key string, value int) error
}

func newLimitsStorage() limitsStorage {
	return newInMemoryLimitsStorage()
}
