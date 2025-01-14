package limiter

import (
	"errors"
)

var NoAvailableTokenError = errors.New("no available token")
var LimiterClosedError = errors.New("limiter closed")
