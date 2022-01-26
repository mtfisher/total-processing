package core

import "time"

func (c Core) GetCacheValue(key string) string {
	val, found := c.c.Get(key)
	if found {
		return val.(string)
	}

	return ""
}

func (c Core) SetCacheValue(key string, val interface{}, ttl time.Duration) {
	c.c.Set(key, val, ttl)
}

func (c Core) RemoveCacheKey(key string) {
	c.c.Delete(key)
}
