package windivert

import "syscall"

// Common errors
var (
	ErrNoData          = Error(syscall.ERROR_NO_DATA)
	ErrHostUnreachable = Error(syscall.ERROR_HOST_UNREACHABLE)
)
