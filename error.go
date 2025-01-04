package windivert

import (
	"syscall"
)

const (
	ERROR_NO_DATA          = syscall.Errno(232)  // ERROR_NO_DATA
	ERROR_HOST_UNREACHABLE = syscall.Errno(1232) // ERROR_HOST_UNREACHABLE
)

// Error represents a WinDivert error
type Error syscall.Errno

// Error returns the error string
func (e Error) Error() string {
	return syscall.Errno(e).Error()
}

// Common errors
var (
	ErrNoData          = Error(ERROR_NO_DATA)
	ErrHostUnreachable = Error(ERROR_HOST_UNREACHABLE)
)
