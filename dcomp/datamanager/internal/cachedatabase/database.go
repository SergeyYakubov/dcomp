// Package containes a database objects and functions to work with it.
// db is an interface to a specific implementation (currently implemented mongodb and mockdatabase used for tests)
package cachedatabase

import (
	"github.com/sergeyyakubov/dcomp/dcomp/server"
	"time"
)

type CachedFile struct {
	SHA      string
	Size     int
	LastUsed time.Time
}

type Agent interface {
	AddFile(file CachedFile) error
	FileInCache(sha string) (bool, error)
	CleanCache(unusedSince time.Duration) ([]string, error)
	GetCacheSize() (int, error)
	Connect() error
	SetDefaults(...interface{})
	SetServer(*server.Server)
	Close() error
}
