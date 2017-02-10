// +build !release

package cachedatabase

import (
	"time"

	"github.com/sergeyyakubov/dcomp/dcomp/server"
)

type Mockdatabase struct {
}

func (db *Mockdatabase) AddFile(file CachedFile) error {
	return nil
}

func (db *Mockdatabase) FileInCache(sha string) (bool, error) {
	return true, nil
}

func (db *Mockdatabase) CleanCache(unusedSince time.Duration) (filesToDelete []string, err error) {
	return
}

func (db *Mockdatabase) GetCacheSize() (int, error) {
	return 0, nil
}

func (db *Mockdatabase) Connect() error {
	return nil
}

func (db *Mockdatabase) SetDefaults(...interface{}) {

}
func (db *Mockdatabase) SetServer(*server.Server) {

}
func (db *Mockdatabase) Close() error {
	return nil
}
