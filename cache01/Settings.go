package cache01

import "time"

type Settings struct {
	// disables the cache. cache files, though, should not be deleted until
	// CacheSingleFileLifetime passes
	PassThrough bool

	ListDirTimeout time.Duration
}

func MakeDefaultSettings() *Settings {
	return &Settings{
		ListDirTimeout: time.Duration(3 * (time.Hour * 24)),
	}
}
