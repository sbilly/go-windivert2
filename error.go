package windivert

/*
#include <windows.h>
#include <windivert.h>

HANDLE getLastError() {
    return GetLastError();
}
*/
import "C"

import (
	"fmt"
	"syscall"
)

// getLastError returns the last error that occurred
func getLastError() error {
	code := C.getLastError()
	if code == 0 {
		return nil
	}
	return fmt.Errorf("windivert error: %d", uint32(code))
}

var (
	ErrNoData          = syscall.WSAEWOULDBLOCK
	ErrHostUnreachable = syscall.WSAEHOSTUNREACH
)
