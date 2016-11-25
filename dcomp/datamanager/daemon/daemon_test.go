package daemon

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetDaemonConfiguration(t *testing.T) {

	configFileName := `daemon_test.yaml`

	err := setDaemonConfiguration(configFileName)

	assert.Nil(t, err, "Should not be error")

	assert.Equal(t, settings.Daemon.Addr, ":8006", "addr")
	assert.Equal(t, settings.Daemon.Key, "12345", "key")

	assert.Equal(t, settings.Dcompd.Host, "localhost", "host")
	assert.Equal(t, settings.Dcompd.Port, 8001, "port")

	assert.Equal(t, settings.Resource.BaseDir, "/home/yakubov/dcomp_test", "basedir")

	configFileName = `aaa`

	err = setDaemonConfiguration(configFileName)
	assert.NotNil(t, err, "Should be error")

	configFileName = `/etc/dcomp/plugins/local/local_dmd.yaml`
	setDaemonConfiguration(configFileName)
}
