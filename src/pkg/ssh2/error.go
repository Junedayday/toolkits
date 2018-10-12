package ssh2

import "errors"

var (
	errNoAuth         = errors.New("No password or private key")
	errClientTimeout  = errors.New("client connect time out")
	errSessionTimeout = errors.New("new session time out")
	errCmdTimeout     = errors.New("run cmd time out")
)
