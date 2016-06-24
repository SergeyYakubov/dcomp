package cli

import ()

func SetServerConfiguration() error {
	Server.host = "localhost"
	Server.port = 8000
	return nil
}
