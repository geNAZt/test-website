package cache

import "time"

func NewFaviconCache() (*TimeoutCache, error) {
	return NewTimeoutCache(int64(60) * int64(time.Minute))
}
