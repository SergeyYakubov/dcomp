package cli

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSetServerConfiguration(t *testing.T) {
	SetDaemonConfiguration()
	assert.Equal(t, "localhost", daemon.Host, "")
	assert.Equal(t, 8001, daemon.Port, "")
}
