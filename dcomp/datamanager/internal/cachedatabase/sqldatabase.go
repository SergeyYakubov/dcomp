// +build !release

package cachedatabase

import (
	"time"

	"github.com/sergeyyakubov/dcomp/dcomp/server"
)

type SqlDatabase struct {
}

func (db *SqlDatabase) AddFile(file CachedFile) error {
	return nil
}

func (db *SqlDatabase) FileInCache(sha string) (bool, error) {
	return true, nil
}

func (db *SqlDatabase) CleanCache(unusedSince time.Duration) (filesToDelete []string, err error) {
	return
}

func (db *SqlDatabase) GetCacheSize() (int, error) {
	return 0, nil
}

func (db *SqlDatabase) Connect() error {
	return nil
}

func (db *SqlDatabase) SetDefaults(...interface{}) {

}
func (db *SqlDatabase) SetServer(*server.Server) {

}
func (db *SqlDatabase) Close() error {
	return nil
}
