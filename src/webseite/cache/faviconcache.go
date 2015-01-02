package cache

import "time"

func NewFaviconCache() (*TimeoutCache, error) {
	return NewTimeoutCache(int64(60) * time.Minute)
}
