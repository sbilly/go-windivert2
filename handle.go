package windivert

/*
#include <windows.h>
#include <windivert.h>
*/
import "C"

// Handle 的方法实现...
func (h *Handle) Lock() {
	h.mutex.Lock()
}

// 其他方法保持不变...
