package cli

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSetServerConfiguration(t *testing.T) {
	SetServerConfiguration()
	assert.Equal(t, "localhost", Server.Host, "")
	assert.Equal(t, 8000, Server.Port, "")
}
