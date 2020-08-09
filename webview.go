// +build windows

package webview2

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"sync/atomic"
	"syscall"
	"unicode/utf16"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	ole32               = windows.NewLazySystemDLL("ole32")
	ole32CoInitializeEx = ole32.NewProc("CoInitializeEx")

	kernel32                   = windows.NewLazySystemDLL("kernel32")
	kernel32GetProcessHeap     = kernel32.NewProc("GetProcessHeap")
	kernel32HeapAlloc          = kernel32.NewProc("HeapAlloc")
	kernel32HeapFree           = kernel32.NewProc("HeapFree")
	kernel32GetCurrentThreadID = kernel32.NewProc("GetCurrentThreadId")

	user32                   = windows.NewLazySystemDLL("user32")
	user32LoadImageW         = user32.NewProc("LoadImageW")
	user32GetSystemMetrics   = user32.NewProc("GetSystemMetrics")
	user32RegisterClassExW   = user32.NewProc("RegisterClassExW")
	user32CreateWindowExW    = user32.NewProc("CreateWindowExW")
	user32DestroyWindow      = user32.NewProc("DestroyWindow")
	user32ShowWindow         = user32.NewProc("ShowWindow")
	user32UpdateWindow       = user32.NewProc("UpdateWindow")
	user32SetFocus           = user32.NewProc("SetFocus")
	user32GetMessageW        = user32.NewProc("GetMessageW")
	user32TranslateMessage   = user32.NewProc("TranslateMessage")
	user32DispatchMessageW   = user32.NewProc("DispatchMessageW")
	user32DefWindowProcW     = user32.NewProc("DefWindowProcW")
	user32GetClientRect      = user32.NewProc("GetClientRect")
	user32PostQuitMessage    = user32.NewProc("PostQuitMessage")
	user32SetWindowTextW     = user32.NewProc("SetWindowTextW")
	user32PostThreadMessageW = user32.NewProc("PostThreadMessageW")

	defaultHeap uintptr
)

var (
	windowContext     = map[uintptr]interface{}{}
	windowContextSync sync.RWMutex
)

func getWindowContext(wnd uintptr) interface{} {
	windowContextSync.RLock()
	defer windowContextSync.RUnlock()
	return windowContext[wnd]
}

func setWindowContext(wnd uintptr, data interface{}) {
	windowContextSync.Lock()
	defer windowContextSync.Unlock()
	windowContext[wnd] = data
}

const (
	_SystemMetricsCxIcon = 11
	_SystemMetricsCyIcon = 12
)

const (
	_SWShow = 5
)

const (
	_WMDestroy       = 0x0002
	_WMSize          = 0x0005
	_WMClose         = 0x0010
	_WMQuit          = 0x0012
	_WMGetMinMaxInfo = 0x0024
	_WMApp           = 0x8000
)

type _WndClassExW struct {
	cbSize        uint32
	style         uint32
	lpfnWndProc   uintptr
	cnClsExtra    int32
	cbWndExtra    int32
	hInstance     windows.Handle
	hIcon         windows.Handle
	hCursor       windows.Handle
	hbrBackground windows.Handle
	lpszMenuName  *uint16
	lpszClassName *uint16
	hIconSm       windows.Handle
}

type _Rect struct {
	Left   int32
	Top    int32
	Right  int32
	Bottom int32
}

type _Point struct {
	x, y int32
}

type _Msg struct {
	hwnd     syscall.Handle
	message  uint32
	wParam   uintptr
	lParam   uintptr
	time     uint32
	pt       _Point
	lPrivate uint32
}

func init() {
	runtime.LockOSThread()

	r, _, _ := ole32CoInitializeEx.Call(0, 2)
	if r < 0 {
		log.Printf("Warning: CoInitializeEx call failed: E=%08x", r)
	}

	defaultHeap, _, _ = kernel32GetProcessHeap.Call()
}

func utf16PtrToString(p *uint16) string {
	if p == nil {
		return ""
	}
	// Find NUL terminator.
	end := unsafe.Pointer(p)
	n := 0
	for *(*uint16)(end) != 0 {
		end = unsafe.Pointer(uintptr(end) + unsafe.Sizeof(*p))
		n++
	}
	s := (*[(1 << 30) - 1]uint16)(unsafe.Pointer(p))[:n:n]
	return string(utf16.Decode(s))
}

type chromiumedge struct {
	hwnd                uintptr
	controller          *iCoreWebView2Controller
	webview             *iCoreWebView2
	inited              uintptr
	envCompleted        *iCoreWebView2CreateCoreWebView2EnvironmentCompletedHandler
	controllerCompleted *iCoreWebView2CreateCoreWebView2ControllerCompletedHandler
	webMessageReceived  *iCoreWebView2WebMessageReceivedEventHandler
	permissionRequested *iCoreWebView2PermissionRequestedEventHandler
	msgcb               func(string)
}

type browser interface {
	Embed(hwnd uintptr) bool
	Resize()
	Navigate(url string)
	Init(script string)
	Eval(script string)
}

type webview struct {
	hwnd       uintptr
	mainthread uintptr
	browser    browser
}

func newchromiumedge() *chromiumedge {
	e := &chromiumedge{}
	e.envCompleted = newICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler(e)
	e.controllerCompleted = newICoreWebView2CreateCoreWebView2ControllerCompletedHandler(e)
	e.webMessageReceived = newICoreWebView2WebMessageReceivedEventHandler(e)
	e.permissionRequested = newICoreWebView2PermissionRequestedEventHandler(e)
	return e
}

func (e *chromiumedge) Embed(hwnd uintptr) bool {
	e.hwnd = hwnd
	currentExePath := make([]uint16, windows.MAX_PATH)
	windows.GetModuleFileName(windows.Handle(0), &currentExePath[0], windows.MAX_PATH)
	currentExeName := filepath.Base(windows.UTF16ToString(currentExePath))
	dataPath := filepath.Join(os.Getenv("AppData"), currentExeName)
	res := createCoreWebView2EnvironmentWithOptions(nil, windows.StringToUTF16Ptr(dataPath), 0, e.envCompleted)
	if res != 0 {
		return false
	}
	var msg _Msg
	for {
		if atomic.LoadUintptr(&e.inited) != 0 {
			break
		}
		r, _, _ := user32GetMessageW.Call(
			uintptr(unsafe.Pointer(&msg)),
			0,
			0,
			0,
		)
		if r == 0 {
			break
		}
		user32TranslateMessage.Call(uintptr(unsafe.Pointer(&msg)))
		user32DispatchMessageW.Call(uintptr(unsafe.Pointer(&msg)))
	}
	e.Init("window.external={invoke:s=>window.chrome.webview.postMessage(s)}")
	return true
}

func (e *chromiumedge) Resize() {
	if e.controller == nil {
		return
	}
	var bounds _Rect
	user32GetClientRect.Call(e.hwnd, uintptr(unsafe.Pointer(&bounds)))
	e.controller.vtbl.PutBounds.Call(
		uintptr(unsafe.Pointer(e.controller)),
		uintptr(unsafe.Pointer(&bounds)),
	)
}

func (e *chromiumedge) Navigate(url string) {
	e.webview.vtbl.Navigate.Call(
		uintptr(unsafe.Pointer(e.webview)),
		uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(url))),
	)
}

func (e *chromiumedge) Init(script string) {
	e.webview.vtbl.AddScriptToExecuteOnDocumentCreated.Call(
		uintptr(unsafe.Pointer(e.webview)),
		uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(script))),
		0,
	)
}

func (e *chromiumedge) Eval(script string) {
	e.webview.vtbl.ExecuteScript.Call(
		uintptr(unsafe.Pointer(e.webview)),
		uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(script))),
		0,
	)
}

func (e *chromiumedge) QueryInterface(refiid, object uintptr) uintptr {
	return 0
}

func (e *chromiumedge) AddRef() uintptr {
	return 1
}

func (e *chromiumedge) Release() uintptr {
	return 1
}

func (e *chromiumedge) EnvironmentCompleted(res uintptr, env *iCoreWebView2Environment) uintptr {
	env.vtbl.CreateCoreWebView2Controller.Call(
		uintptr(unsafe.Pointer(env)),
		e.hwnd,
		uintptr(unsafe.Pointer(e.controllerCompleted)),
	)
	return 0
}

func (e *chromiumedge) ControllerCompleted(res uintptr, controller *iCoreWebView2Controller) uintptr {
	controller.vtbl.AddRef.Call(uintptr(unsafe.Pointer(controller)))

	var webview *iCoreWebView2
	var token _EventRegistrationToken
	controller.vtbl.GetCoreWebView2.Call(
		uintptr(unsafe.Pointer(controller)),
		uintptr(unsafe.Pointer(&webview)),
	)
	webview.vtbl.AddWebMessageReceived.Call(
		uintptr(unsafe.Pointer(webview)),
		uintptr(unsafe.Pointer(e.webMessageReceived)),
		uintptr(unsafe.Pointer(&token)),
	)
	webview.vtbl.AddPermissionRequested.Call(
		uintptr(unsafe.Pointer(webview)),
		uintptr(unsafe.Pointer(e.permissionRequested)),
		uintptr(unsafe.Pointer(&token)),
	)

	e.controller = controller
	e.controller.vtbl.GetCoreWebView2.Call(
		uintptr(unsafe.Pointer(e.controller)),
		uintptr(unsafe.Pointer(&e.webview)),
	)
	e.webview.vtbl.AddRef.Call(
		uintptr(unsafe.Pointer(e.webview)),
	)

	atomic.StoreUintptr(&e.inited, 1)

	return 0
}

func (e *chromiumedge) MessageReceived(sender *iCoreWebView2, args *iCoreWebView2WebMessageReceivedEventArgs) uintptr {
	var message *uint16
	args.vtbl.TryGetWebMessageAsString.Call(
		uintptr(unsafe.Pointer(args)),
		uintptr(unsafe.Pointer(&message)),
	)
	e.msgcb(utf16PtrToString(message))
	sender.vtbl.PostWebMessageAsString.Call(
		uintptr(unsafe.Pointer(sender)),
		uintptr(unsafe.Pointer(message)),
	)
	windows.CoTaskMemFree(unsafe.Pointer(message))
	return 0
}

func (e *chromiumedge) PermissionRequested(sender *iCoreWebView2, args *iCoreWebView2PermissionRequestedEventArgs) uintptr {
	var kind _CoreWebView2PermissionKind
	args.vtbl.GetPermissionKind.Call(
		uintptr(unsafe.Pointer(args)),
		uintptr(kind),
	)
	if kind == _CoreWebView2PermissionKindClipboardRead {
		args.vtbl.PutState.Call(
			uintptr(unsafe.Pointer(args)),
			uintptr(_CoreWebView2PermissionStateAllow),
		)
	}
	return 0
}

// New creates a new webview in a new window.
func New(debug bool) WebView { return NewWindow(debug, nil) }

// NewWindow creates a new webview using an existing window.
func NewWindow(debug bool, window unsafe.Pointer) WebView {
	w := &webview{}
	w.browser = newchromiumedge()
	w.mainthread, _, _ = kernel32GetCurrentThreadID.Call()
	if !w.Create(debug, window) {
		return nil
	}
	return w
}

func wndproc(hwnd, msg, wp, lp uintptr) uintptr {
	if w, ok := getWindowContext(hwnd).(*webview); ok {
		switch msg {
		case _WMSize:
			w.browser.Resize()
		case _WMClose:
			user32DestroyWindow.Call(hwnd)
		case _WMDestroy:
			w.Terminate()
		case _WMGetMinMaxInfo:
			// TODO
		default:
			r, _, _ := user32DefWindowProcW.Call(hwnd, msg, wp, lp)
			return r
		}
		return 0
	}
	r, _, _ := user32DefWindowProcW.Call(hwnd, msg, wp, lp)
	return r
}

func (w *webview) Create(debug bool, window unsafe.Pointer) bool {
	var hinstance windows.Handle
	windows.GetModuleHandleEx(0, nil, &hinstance)

	icow, _, _ := user32GetSystemMetrics.Call(_SystemMetricsCxIcon)
	icoh, _, _ := user32GetSystemMetrics.Call(_SystemMetricsCyIcon)

	icon, _, _ := user32LoadImageW.Call(uintptr(hinstance), 32512, icow, icoh, 0)

	wc := _WndClassExW{
		cbSize:        uint32(unsafe.Sizeof(_WndClassExW{})),
		hInstance:     hinstance,
		lpszClassName: windows.StringToUTF16Ptr("webview"),
		hIcon:         windows.Handle(icon),
		hIconSm:       windows.Handle(icon),
		lpfnWndProc:   windows.NewCallback(wndproc),
	}
	user32RegisterClassExW.Call(uintptr(unsafe.Pointer(&wc)))
	w.hwnd, _, _ = user32CreateWindowExW.Call(
		0,
		uintptr(unsafe.Pointer(windows.StringToUTF16Ptr("webview"))),
		uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(""))),
		0xCF0000,   // WS_OVERLAPPEDWINDOW
		0x80000000, // CW_USEDEFAULT
		0x80000000, // CW_USEDEFAULT
		640,
		480,
		0,
		0,
		uintptr(hinstance),
		0,
	)
	setWindowContext(w.hwnd, w)

	user32ShowWindow.Call(w.hwnd, _SWShow)
	user32UpdateWindow.Call(w.hwnd)
	user32SetFocus.Call(w.hwnd)

	if !w.browser.Embed(w.hwnd) {
		return false
	}
	w.browser.Resize()
	return true
}

func (w *webview) Destroy() {
}

func (w *webview) Run() {
	var msg _Msg
	for {
		user32GetMessageW.Call(
			uintptr(unsafe.Pointer(&msg)),
			0,
			0,
			0,
		)
		if msg.message == _WMApp {

		} else if msg.message == _WMQuit {
			return
		}
		user32TranslateMessage.Call(uintptr(unsafe.Pointer(&msg)))
		user32DispatchMessageW.Call(uintptr(unsafe.Pointer(&msg)))
	}
}

func (w *webview) Terminate() {
	user32PostQuitMessage.Call(0)
}

func (w *webview) Window() unsafe.Pointer {
	return unsafe.Pointer(w.hwnd)
}

func (w *webview) Navigate(url string) {
	w.browser.Navigate(url)
}

func (w *webview) SetTitle(title string) {
	user32SetWindowTextW.Call(w.hwnd, uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(title))))
}

func (w *webview) SetSize(width int, height int, hint Hint) {
	// TODO
}

func (w *webview) Init(js string) {
	w.browser.Init(js)
}

func (w *webview) Eval(js string) {
	w.browser.Eval(js)
}

func (w *webview) Dispatch(f func()) {
	// TODO
}

func (w *webview) Bind(name string, f interface{}) error {
	// TODO
	return nil
}
