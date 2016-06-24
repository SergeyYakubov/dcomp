package cli

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSetServerConfiguration(t *testing.T) {
	SetServerConfiguration()
	assert.Equal(t, "localhost", Server.host, "")
	assert.Equal(t, 8000, Server.port, "")
}
