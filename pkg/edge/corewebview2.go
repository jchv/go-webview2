// +build windows

package edge

import (
	"github.com/jchv/go-webview2/internal/w32"
	"log"
	"runtime"
	"syscall"
	"unsafe"

	"github.com/jchv/go-webview2/webviewloader"
	"golang.org/x/sys/windows"
)

func init() {
	runtime.LockOSThread()

	r, _, _ := w32.Ole32CoInitializeEx.Call(0, 2)
	if int(r) < 0 {
		log.Printf("Warning: CoInitializeEx call failed: E=%08x", r)
	}
}

type _EventRegistrationToken struct {
	value int64
}

type _CoreWebView2PermissionKind uint32

const (
	_CoreWebView2PermissionKindUnknownPermission _CoreWebView2PermissionKind = iota
	_CoreWebView2PermissionKindMicrophone
	_CoreWebView2PermissionKindCamera
	_CoreWebView2PermissionKindGeolocation
	_CoreWebView2PermissionKindNotifications
	_CoreWebView2PermissionKindOtherSensors
	_CoreWebView2PermissionKindClipboardRead
)

type _CoreWebView2PermissionState uint32

const (
	_CoreWebView2PermissionStateDefault _CoreWebView2PermissionState = iota
	_CoreWebView2PermissionStateAllow
	_CoreWebView2PermissionStateDeny
)

func createCoreWebView2EnvironmentWithOptions(browserExecutableFolder, userDataFolder *uint16, environmentOptions uintptr, environmentCompletedHandle *iCoreWebView2CreateCoreWebView2EnvironmentCompletedHandler) (uintptr, error) {
	return webviewloader.CreateCoreWebView2EnvironmentWithOptions(
		browserExecutableFolder,
		userDataFolder,
		environmentOptions,
		uintptr(unsafe.Pointer(environmentCompletedHandle)),
	)
}

// ComProc stores a COM procedure.
type ComProc uintptr

// NewComProc creates a new COM proc from a Go function.
func NewComProc(fn interface{}) ComProc {
	return ComProc(windows.NewCallback(fn))
}

// Call calls a COM procedure.
func (p ComProc) Call(a ...uintptr) (r1, r2 uintptr, lastErr error) {
	switch len(a) {
	case 0:
		return syscall.Syscall(uintptr(p), 0, 0, 0, 0)
	case 1:
		return syscall.Syscall(uintptr(p), 1, a[0], 0, 0)
	case 2:
		return syscall.Syscall(uintptr(p), 2, a[0], a[1], 0)
	case 3:
		return syscall.Syscall(uintptr(p), 3, a[0], a[1], a[2])
	case 4:
		return syscall.Syscall6(uintptr(p), 4, a[0], a[1], a[2], a[3], 0, 0)
	case 5:
		return syscall.Syscall6(uintptr(p), 5, a[0], a[1], a[2], a[3], a[4], 0)
	case 6:
		return syscall.Syscall6(uintptr(p), 6, a[0], a[1], a[2], a[3], a[4], a[5])
	case 7:
		return syscall.Syscall9(uintptr(p), 7, a[0], a[1], a[2], a[3], a[4], a[5], a[6], 0, 0)
	case 8:
		return syscall.Syscall9(uintptr(p), 8, a[0], a[1], a[2], a[3], a[4], a[5], a[6], a[7], 0)
	case 9:
		return syscall.Syscall9(uintptr(p), 9, a[0], a[1], a[2], a[3], a[4], a[5], a[6], a[7], a[8])
	case 10:
		return syscall.Syscall12(uintptr(p), 10, a[0], a[1], a[2], a[3], a[4], a[5], a[6], a[7], a[8], a[9], 0, 0)
	case 11:
		return syscall.Syscall12(uintptr(p), 11, a[0], a[1], a[2], a[3], a[4], a[5], a[6], a[7], a[8], a[9], a[10], 0)
	case 12:
		return syscall.Syscall12(uintptr(p), 12, a[0], a[1], a[2], a[3], a[4], a[5], a[6], a[7], a[8], a[9], a[10], a[11])
	case 13:
		return syscall.Syscall15(uintptr(p), 13, a[0], a[1], a[2], a[3], a[4], a[5], a[6], a[7], a[8], a[9], a[10], a[11], a[12], 0, 0)
	case 14:
		return syscall.Syscall15(uintptr(p), 14, a[0], a[1], a[2], a[3], a[4], a[5], a[6], a[7], a[8], a[9], a[10], a[11], a[12], a[13], 0)
	case 15:
		return syscall.Syscall15(uintptr(p), 15, a[0], a[1], a[2], a[3], a[4], a[5], a[6], a[7], a[8], a[9], a[10], a[11], a[12], a[13], a[14])
	default:
		panic("too many arguments")
	}
}

// IUnknown

type iUnknownVtbl struct {
	QueryInterface ComProc
	AddRef         ComProc
	Release        ComProc
}

type iUnknownImpl interface {
	QueryInterface(refiid, object uintptr) uintptr
	AddRef() uintptr
	Release() uintptr
}

// ICoreWebView2

type iCoreWebView2Vtbl struct {
	iUnknownVtbl
	GetSettings                            ComProc
	GetSource                              ComProc
	Navigate                               ComProc
	NavigateToString                       ComProc
	AddNavigationStarting                  ComProc
	RemoveNavigationStarting               ComProc
	AddContentLoading                      ComProc
	RemoveContentLoading                   ComProc
	AddSourceChanged                       ComProc
	RemoveSourceChanged                    ComProc
	AddHistoryChanged                      ComProc
	RemoveHistoryChanged                   ComProc
	AddNavigationCompleted                 ComProc
	RemoveNavigationCompleted              ComProc
	AddFrameNavigationStarting             ComProc
	RemoveFrameNavigationStarting          ComProc
	AddFrameNavigationCompleted            ComProc
	RemoveFrameNavigationCompleted         ComProc
	AddScriptDialogOpening                 ComProc
	RemoveScriptDialogOpening              ComProc
	AddPermissionRequested                 ComProc
	RemovePermissionRequested              ComProc
	AddProcessFailed                       ComProc
	RemoveProcessFailed                    ComProc
	AddScriptToExecuteOnDocumentCreated    ComProc
	RemoveScriptToExecuteOnDocumentCreated ComProc
	ExecuteScript                          ComProc
	CapturePreview                         ComProc
	Reload                                 ComProc
	PostWebMessageAsJSON                   ComProc
	PostWebMessageAsString                 ComProc
	AddWebMessageReceived                  ComProc
	RemoveWebMessageReceived               ComProc
	CallDevToolsProtocolMethod             ComProc
	GetBrowserProcessID                    ComProc
	GetCanGoBack                           ComProc
	GetCanGoForward                        ComProc
	GoBack                                 ComProc
	GoForward                              ComProc
	GetDevToolsProtocolEventReceiver       ComProc
	Stop                                   ComProc
	AddNewWindowRequested                  ComProc
	RemoveNewWindowRequested               ComProc
	AddDocumentTitleChanged                ComProc
	RemoveDocumentTitleChanged             ComProc
	GetDocumentTitle                       ComProc
	AddHostObjectToScript                  ComProc
	RemoveHostObjectFromScript             ComProc
	OpenDevToolsWindow                     ComProc
	AddContainsFullScreenElementChanged    ComProc
	RemoveContainsFullScreenElementChanged ComProc
	GetContainsFullScreenElement           ComProc
	AddWebResourceRequested                ComProc
	RemoveWebResourceRequested             ComProc
	AddWebResourceRequestedFilter          ComProc
	RemoveWebResourceRequestedFilter       ComProc
	AddWindowCloseRequested                ComProc
	RemoveWindowCloseRequested             ComProc
}

type iCoreWebView2 struct {
	vtbl *iCoreWebView2Vtbl
}

// ICoreWebView2Environment

type iCoreWebView2EnvironmentVtbl struct {
	iUnknownVtbl
	CreateCoreWebView2Controller     ComProc
	CreateWebResourceResponse        ComProc
	GetBrowserVersionString          ComProc
	AddNewBrowserVersionAvailable    ComProc
	RemoveNewBrowserVersionAvailable ComProc
}

type iCoreWebView2Environment struct {
	vtbl *iCoreWebView2EnvironmentVtbl
}

// ICoreWebView2Controller

type iCoreWebView2ControllerVtbl struct {
	iUnknownVtbl
	GetIsVisible                      ComProc
	PutIsVisible                      ComProc
	GetBounds                         ComProc
	PutBounds                         ComProc
	GetZoomFactor                     ComProc
	PutZoomFactor                     ComProc
	AddZoomFactorChanged              ComProc
	RemoveZoomFactorChanged           ComProc
	SetBoundsAndZoomFactor            ComProc
	MoveFocus                         ComProc
	AddMoveFocusRequested             ComProc
	RemoveMoveFocusRequested          ComProc
	AddGotFocus                       ComProc
	RemoveGotFocus                    ComProc
	AddLostFocus                      ComProc
	RemoveLostFocus                   ComProc
	AddAcceleratorKeyPressed          ComProc
	RemoveAcceleratorKeyPressed       ComProc
	GetParentWindow                   ComProc
	PutParentWindow                   ComProc
	NotifyParentWindowPositionChanged ComProc
	Close                             ComProc
	GetCoreWebView2                   ComProc
}

type iCoreWebView2Controller struct {
	vtbl *iCoreWebView2ControllerVtbl
}

// ICoreWebView2WebMessageReceivedEventArgs

type iCoreWebView2WebMessageReceivedEventArgsVtbl struct {
	iUnknownVtbl
	GetSource                ComProc
	GetWebMessageAsJSON      ComProc
	TryGetWebMessageAsString ComProc
}

type iCoreWebView2WebMessageReceivedEventArgs struct {
	vtbl *iCoreWebView2WebMessageReceivedEventArgsVtbl
}

// ICoreWebView2PermissionRequestedEventArgs

type iCoreWebView2PermissionRequestedEventArgsVtbl struct {
	iUnknownVtbl
	GetURI             ComProc
	GetPermissionKind  ComProc
	GetIsUserInitiated ComProc
	GetState           ComProc
	PutState           ComProc
	GetDeferral        ComProc
}

type iCoreWebView2PermissionRequestedEventArgs struct {
	vtbl *iCoreWebView2PermissionRequestedEventArgsVtbl
}

// ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler

type iCoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerImpl interface {
	iUnknownImpl
	EnvironmentCompleted(res uintptr, env *iCoreWebView2Environment) uintptr
}

type iCoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerVtbl struct {
	iUnknownVtbl
	Invoke ComProc
}

type iCoreWebView2CreateCoreWebView2EnvironmentCompletedHandler struct {
	vtbl *iCoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerVtbl
	impl iCoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerImpl
}

func _ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerIUnknownQueryInterface(this *iCoreWebView2CreateCoreWebView2EnvironmentCompletedHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func _ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerIUnknownAddRef(this *iCoreWebView2CreateCoreWebView2EnvironmentCompletedHandler) uintptr {
	return this.impl.AddRef()
}

func _ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerIUnknownRelease(this *iCoreWebView2CreateCoreWebView2EnvironmentCompletedHandler) uintptr {
	return this.impl.Release()
}

func _ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerInvoke(this *iCoreWebView2CreateCoreWebView2EnvironmentCompletedHandler, res uintptr, env *iCoreWebView2Environment) uintptr {
	return this.impl.EnvironmentCompleted(res, env)
}

var iCoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerFn = iCoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerVtbl{
	iUnknownVtbl{
		NewComProc(_ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerIUnknownQueryInterface),
		NewComProc(_ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerIUnknownAddRef),
		NewComProc(_ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerIUnknownRelease),
	},
	NewComProc(_ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerInvoke),
}

func newICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler(impl iCoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerImpl) *iCoreWebView2CreateCoreWebView2EnvironmentCompletedHandler {
	return &iCoreWebView2CreateCoreWebView2EnvironmentCompletedHandler{
		vtbl: &iCoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerFn,
		impl: impl,
	}
}

// ICoreWebView2CreateCoreWebView2ControllerCompletedHandler

type iCoreWebView2CreateCoreWebView2ControllerCompletedHandlerImpl interface {
	iUnknownImpl
	ControllerCompleted(res uintptr, controller *iCoreWebView2Controller) uintptr
}

type iCoreWebView2CreateCoreWebView2ControllerCompletedHandlerVtbl struct {
	iUnknownVtbl
	Invoke ComProc
}

type iCoreWebView2CreateCoreWebView2ControllerCompletedHandler struct {
	vtbl *iCoreWebView2CreateCoreWebView2ControllerCompletedHandlerVtbl
	impl iCoreWebView2CreateCoreWebView2ControllerCompletedHandlerImpl
}

func _ICoreWebView2CreateCoreWebView2ControllerCompletedHandlerIUnknownQueryInterface(this *iCoreWebView2CreateCoreWebView2ControllerCompletedHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func _ICoreWebView2CreateCoreWebView2ControllerCompletedHandlerIUnknownAddRef(this *iCoreWebView2CreateCoreWebView2ControllerCompletedHandler) uintptr {
	return this.impl.AddRef()
}

func _ICoreWebView2CreateCoreWebView2ControllerCompletedHandlerIUnknownRelease(this *iCoreWebView2CreateCoreWebView2ControllerCompletedHandler) uintptr {
	return this.impl.Release()
}

func _ICoreWebView2CreateCoreWebView2ControllerCompletedHandlerInvoke(this *iCoreWebView2CreateCoreWebView2ControllerCompletedHandler, res uintptr, controller *iCoreWebView2Controller) uintptr {
	return this.impl.ControllerCompleted(res, controller)
}

var iCoreWebView2CreateCoreWebView2ControllerCompletedHandlerFn = iCoreWebView2CreateCoreWebView2ControllerCompletedHandlerVtbl{
	iUnknownVtbl{
		NewComProc(_ICoreWebView2CreateCoreWebView2ControllerCompletedHandlerIUnknownQueryInterface),
		NewComProc(_ICoreWebView2CreateCoreWebView2ControllerCompletedHandlerIUnknownAddRef),
		NewComProc(_ICoreWebView2CreateCoreWebView2ControllerCompletedHandlerIUnknownRelease),
	},
	NewComProc(_ICoreWebView2CreateCoreWebView2ControllerCompletedHandlerInvoke),
}

func newICoreWebView2CreateCoreWebView2ControllerCompletedHandler(impl iCoreWebView2CreateCoreWebView2ControllerCompletedHandlerImpl) *iCoreWebView2CreateCoreWebView2ControllerCompletedHandler {
	return &iCoreWebView2CreateCoreWebView2ControllerCompletedHandler{
		vtbl: &iCoreWebView2CreateCoreWebView2ControllerCompletedHandlerFn,
		impl: impl,
	}
}

// ICoreWebView2WebMessageReceivedEventHandler

type iCoreWebView2WebMessageReceivedEventHandlerImpl interface {
	iUnknownImpl
	MessageReceived(sender *iCoreWebView2, args *iCoreWebView2WebMessageReceivedEventArgs) uintptr
}

type iCoreWebView2WebMessageReceivedEventHandlerVtbl struct {
	iUnknownVtbl
	Invoke ComProc
}

type iCoreWebView2WebMessageReceivedEventHandler struct {
	vtbl *iCoreWebView2WebMessageReceivedEventHandlerVtbl
	impl iCoreWebView2WebMessageReceivedEventHandlerImpl
}

func _ICoreWebView2WebMessageReceivedEventHandlerIUnknownQueryInterface(this *iCoreWebView2WebMessageReceivedEventHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func _ICoreWebView2WebMessageReceivedEventHandlerIUnknownAddRef(this *iCoreWebView2WebMessageReceivedEventHandler) uintptr {
	return this.impl.AddRef()
}

func _ICoreWebView2WebMessageReceivedEventHandlerIUnknownRelease(this *iCoreWebView2WebMessageReceivedEventHandler) uintptr {
	return this.impl.Release()
}

func _ICoreWebView2WebMessageReceivedEventHandlerInvoke(this *iCoreWebView2WebMessageReceivedEventHandler, sender *iCoreWebView2, args *iCoreWebView2WebMessageReceivedEventArgs) uintptr {
	return this.impl.MessageReceived(sender, args)
}

var iCoreWebView2WebMessageReceivedEventHandlerFn = iCoreWebView2WebMessageReceivedEventHandlerVtbl{
	iUnknownVtbl{
		NewComProc(_ICoreWebView2WebMessageReceivedEventHandlerIUnknownQueryInterface),
		NewComProc(_ICoreWebView2WebMessageReceivedEventHandlerIUnknownAddRef),
		NewComProc(_ICoreWebView2WebMessageReceivedEventHandlerIUnknownRelease),
	},
	NewComProc(_ICoreWebView2WebMessageReceivedEventHandlerInvoke),
}

func newICoreWebView2WebMessageReceivedEventHandler(impl iCoreWebView2WebMessageReceivedEventHandlerImpl) *iCoreWebView2WebMessageReceivedEventHandler {
	return &iCoreWebView2WebMessageReceivedEventHandler{
		vtbl: &iCoreWebView2WebMessageReceivedEventHandlerFn,
		impl: impl,
	}
}

// ICoreWebView2PermissionRequestedEventHandler

type iCoreWebView2PermissionRequestedEventHandlerImpl interface {
	iUnknownImpl
	PermissionRequested(sender *iCoreWebView2, args *iCoreWebView2PermissionRequestedEventArgs) uintptr
}

type iCoreWebView2PermissionRequestedEventHandlerVtbl struct {
	iUnknownVtbl
	Invoke ComProc
}

type iCoreWebView2PermissionRequestedEventHandler struct {
	vtbl *iCoreWebView2PermissionRequestedEventHandlerVtbl
	impl iCoreWebView2PermissionRequestedEventHandlerImpl
}

func _ICoreWebView2PermissionRequestedEventHandlerIUnknownQueryInterface(this *iCoreWebView2PermissionRequestedEventHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func _ICoreWebView2PermissionRequestedEventHandlerIUnknownAddRef(this *iCoreWebView2PermissionRequestedEventHandler) uintptr {
	return this.impl.AddRef()
}

func _ICoreWebView2PermissionRequestedEventHandlerIUnknownRelease(this *iCoreWebView2PermissionRequestedEventHandler) uintptr {
	return this.impl.Release()
}

func _ICoreWebView2PermissionRequestedEventHandlerInvoke(this *iCoreWebView2PermissionRequestedEventHandler, sender *iCoreWebView2, args *iCoreWebView2PermissionRequestedEventArgs) uintptr {
	return this.impl.PermissionRequested(sender, args)
}

var iCoreWebView2PermissionRequestedEventHandlerFn = iCoreWebView2PermissionRequestedEventHandlerVtbl{
	iUnknownVtbl{
		NewComProc(_ICoreWebView2PermissionRequestedEventHandlerIUnknownQueryInterface),
		NewComProc(_ICoreWebView2PermissionRequestedEventHandlerIUnknownAddRef),
		NewComProc(_ICoreWebView2PermissionRequestedEventHandlerIUnknownRelease),
	},
	NewComProc(_ICoreWebView2PermissionRequestedEventHandlerInvoke),
}

func newICoreWebView2PermissionRequestedEventHandler(impl iCoreWebView2PermissionRequestedEventHandlerImpl) *iCoreWebView2PermissionRequestedEventHandler {
	return &iCoreWebView2PermissionRequestedEventHandler{
		vtbl: &iCoreWebView2PermissionRequestedEventHandlerFn,
		impl: impl,
	}
}
