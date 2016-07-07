// +build !release

package database

import "errors"

type mockdatabase struct {
}

func (db *mockdatabase) CreateRecord(s interface{}) (string, error) {
	return "1", nil
}
func (db *mockdatabase) Connect(url string) error {
	return nil
}

func (db *mockdatabase) Close() {

}
func (db *mockdatabase) SetDefaults() {
}

func CreateMock() error {
	if db != nil {
		return errors.New("database already created")
	}

	db = new(mockdatabase)

	db.SetDefaults()
	return nil
}
