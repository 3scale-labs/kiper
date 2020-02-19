package threescale

import (
	"time"

	porta "github.com/3scale/3scale-porta-go-client/client"
	gocache "github.com/patrickmn/go-cache"
)

const (
	ttl             = time.Minute
	cleanupInterval = 5 * time.Minute
)

type PortaConfigsCache struct {
	internalStorage *gocache.Cache
}

func newPortaConfigsCache() PortaConfigsCache {
	goCache := gocache.New(ttl, cleanupInterval)
	return PortaConfigsCache{internalStorage: goCache}
}

func (cache *PortaConfigsCache) get(serviceId string, environment string) (*porta.ProxyConfig, bool) {
	config, exists := cache.internalStorage.Get(key(serviceId, environment))

	if !exists {
		return nil, false
	}

	return config.(*porta.ProxyConfig), exists
}

func (cache *PortaConfigsCache) set(serviceId string, environment string, config *porta.ProxyConfig) {
	cache.internalStorage.Set(
		key(serviceId, environment), config, gocache.DefaultExpiration,
	)
}

func key(serviceId string, environment string) string {
	return serviceId + "/" + environment
}
