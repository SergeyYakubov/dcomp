// +build !release

package database

import "errors"

type mockdatabase struct {
}

func (db *mockdatabase) CreateRecord(s interface{}) (string, error) {
	return "1", nil
}
func (db *mockdatabase) Connect(url string) error {
	if url == "localhost:27017" {
		return nil
	}
	return errors.New("mockdb: Server not found")
}

func (db *mockdatabase) Close() {

}
func (db *mockdatabase) SetDefaults() {
}
