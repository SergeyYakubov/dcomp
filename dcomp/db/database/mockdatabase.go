// +build !release

package database

import (
	"errors"

	"reflect"

	"gopkg.in/mgo.v2/bson"
	"stash.desy.de/scm/dc/main.git/dcomp/structs"
)

type mockdatabase struct {
}

func (db *mockdatabase) CreateRecord(s interface{}) (string, error) {
	return "578359205e935a20adb39a18", nil
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

type querryM struct {
	Id bson.ObjectId `bson:"_id"`
}

func (db *mockdatabase) GetRecordByID(id string, records interface{}) error {
	if !bson.IsObjectIdHex(id) {
		return errors.New("wrong id")
	}

	q := querryM{bson.ObjectIdHex(id)}

	return db.GetRecords(&q, records)
}

func (db *mockdatabase) DeleteRecordByID(id string) error {
	if !bson.IsObjectIdHex(id) {
		return errors.New("wrong id")
	}

	if id == "578359205e935a20adb39a18" {
		return nil
	} else {
		return errors.New("job not found")
	}

	return nil
}

func (db *mockdatabase) GetRecords(q interface{}, res interface{}) (err error) {
	data := [...]structs.JobInfo{
		{JobDescription: structs.JobDescription{}, Id: "578359205e935a20adb39a18", Status: 1},
		{JobDescription: structs.JobDescription{}, Id: "578359235e935a21510a2243", Status: 1}}

	resultv := reflect.ValueOf(res)
	if resultv.Kind() != reflect.Ptr || resultv.Elem().Kind() != reflect.Slice {
		panic("result argument must be a slice address")
	}
	slicev := resultv.Elem()
	slicev = slicev.Slice(0, slicev.Cap())
	elemt := slicev.Type().Elem()

	nrec := 0
	for i, v := range data {
		if q == nil || q.(*querryM).Id.Hex() == v.Id {
			nrec++
			if slicev.Len() < nrec {
				elemp := reflect.New(elemt)
				elemp.Elem().Set(reflect.ValueOf(v))
				slicev = reflect.Append(slicev, elemp.Elem())
			} else {
				slicev.Index(i).Elem().Set(reflect.ValueOf(v))
			}

		}
	}
	resultv.Elem().Set(slicev.Slice(0, nrec))
	return nil
}
