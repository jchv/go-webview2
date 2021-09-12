//go:build windows
// +build windows

package edge

import (
	"go-webview2/internal/w32"
	"unsafe"
)

func (e *Chromium) Resize() {
	if e.controller == nil {
		return
	}
	var bounds w32.Rect
	w32.User32GetClientRect.Call(e.hwnd, uintptr(unsafe.Pointer(&bounds)))
	e.controller.vtbl.PutBounds.Call(
		uintptr(unsafe.Pointer(e.controller)),
		uintptr(unsafe.Pointer(&bounds)),
	)
}
