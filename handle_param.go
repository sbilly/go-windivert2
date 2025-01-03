package windivert

// SetParam sets a WinDivert parameter
func (h *Handle) SetParam(param Param, value uint64) error {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if ret := C.WinDivertSetParam(h.handle, C.WINDIVERT_PARAM(param), C.UINT64(value)); ret == 0 {
		return getLastError()
	}
	return nil
}

// GetParam gets a WinDivert parameter
func (h *Handle) GetParam(param Param) (uint64, error) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	var value C.UINT64
	if ret := C.WinDivertGetParam(h.handle, C.WINDIVERT_PARAM(param), &value); ret == 0 {
		return 0, getLastError()
	}
	return uint64(value), nil
}
