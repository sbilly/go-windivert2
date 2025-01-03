package windivert

// Shutdown shuts down a WinDivert handle
func (h *Handle) Shutdown(how Shutdown) error {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if ret := C.WinDivertShutdown(h.handle, C.WINDIVERT_SHUTDOWN(how)); ret == 0 {
		return getLastError()
	}
	return nil
}
