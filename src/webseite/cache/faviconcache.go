package cache

import (
	"time"
	status "github.com/geNAZt/minecraft-status/data"
)

type StoredFavicon struct {
	Favicon  string
	Favicons []status.Favicon
}

var Favicons *TimeoutCache

func init() {
	tempCache, err := NewTimeoutCache(int64(60) * int64(time.Minute))
	if err != nil {
		panic("Could not init favicon cache")
	}

	Favicons = tempCache
}