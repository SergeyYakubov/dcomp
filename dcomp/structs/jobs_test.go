package structs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var setTests = []struct {
	Test     string
	Err      bool
	Messsage string
}{
	{"aaa:bbb", false, "ok"},
	{"aaa:", true, "empty dest"},
	{"aaa:.", true, "relative dest"},
	{"aaa:./", true, "relative dest 2 "},
	{"aaa:/", true, "root dest"},
}

func TestSetTransferFiles(t *testing.T) {

	for _, test := range setTests {
		var setTests TransferFiles
		err := setTests.Set(test.Test)
		if test.Err {
			assert.NotNil(t, err, test.Messsage)
		} else {
			assert.Nil(t, err, test.Messsage)
		}

	}

}
