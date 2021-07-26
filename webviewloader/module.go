package webviewloader

import (
	"fmt"
	"sync"
	"unsafe"

	"github.com/jchv/go-winloader"
	"golang.org/x/sys/windows"
)

var (
	nativeModule = windows.NewLazyDLL("WebView2Loader")
	nativeCreate = nativeModule.NewProc("CreateCoreWebView2EnvironmentWithOptions")

	memOnce   sync.Once
	memModule winloader.Module
	memCreate winloader.Proc
	memErr    error
)

// CreateCoreWebView2EnvironmentWithOptions tries to load WebviewLoader2 and
// call the CreateCoreWebView2EnvironmentWithOptions routine.
func CreateCoreWebView2EnvironmentWithOptions(browserExecutableFolder, userDataFolder *uint16, environmentOptions uintptr, environmentCompletedHandle uintptr) (uintptr, error) {
	nativeErr := nativeModule.Load()
	if nativeErr == nil {
		nativeErr = nativeCreate.Find()
	}
	if nativeErr != nil {
		// DLL is not available natively. Try loading embedded copy.
		memOnce.Do(func() {
			memModule, memErr = winloader.LoadFromMemory(WebView2Loader)
			if memErr == nil {
				memCreate = memModule.Proc("CreateCoreWebView2EnvironmentWithOptions")
			}
		})
		if memErr != nil {
			return 0, fmt.Errorf("Unable to load WebView2Loader.dll from disk: %v -- or from memory: %w", nativeErr, memErr)
		}
		res, _, _ := memCreate.Call(
			uint64(uintptr(unsafe.Pointer(browserExecutableFolder))),
			uint64(uintptr(unsafe.Pointer(userDataFolder))),
			uint64(environmentOptions),
			uint64(environmentCompletedHandle),
		)
		return uintptr(res), nil
	}
	res, _, _ := nativeCreate.Call(
		uintptr(unsafe.Pointer(browserExecutableFolder)),
		uintptr(unsafe.Pointer(userDataFolder)),
		environmentOptions,
		environmentCompletedHandle,
	)
	return res, nil
}
