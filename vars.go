package windivert

import "golang.org/x/sys/windows"

var (
	// WinDivert is the DLL instance
	WinDivert = (*windows.DLL)(nil)
	// WinDivertOpen is the WinDivertOpen procedure
	WinDivertOpen = (*windows.Proc)(nil)
	// WinDivertSys is the path to WinDivert sys file
	WinDivertSys = ""
	// WinDivertDll is the path to WinDivert dll file
	WinDivertDll = ""
	// DeviceName is the WinDivert device name
	DeviceName = windows.StringToUTF16Ptr("WinDivert")
)
