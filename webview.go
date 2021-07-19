// +build windows

package webview2

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
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
	user32GetWindowLongPtrW  = user32.NewProc("GetWindowLongPtrW")
	user32SetWindowLongPtrW  = user32.NewProc("SetWindowLongPtrW")
	user32AdjustWindowRect   = user32.NewProc("AdjustWindowRect")
	user32SetWindowPos       = user32.NewProc("SetWindowPos")
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
	_SWPNoZOrder     = 0x0004
	_SWPNoActivate   = 0x0010
	_SWPNoMove       = 0x0002
	_SWPFrameChanged = 0x0020
)

const (
	_WMDestroy       = 0x0002
	_WMSize          = 0x0005
	_WMClose         = 0x0010
	_WMQuit          = 0x0012
	_WMGetMinMaxInfo = 0x0024
	_WMApp           = 0x8000
)

const (
	_GWLStyle = -16
)

const (
	_WSOverlapped       = 0x00000000
	_WSMaximizeBox      = 0x00020000
	_WSThickFrame       = 0x00040000
	_WSCaption          = 0x00C00000
	_WSSysMenu          = 0x00080000
	_WSMinimizeBox      = 0x00020000
	_WSOverlappedWindow = (_WSOverlapped | _WSCaption | _WSSysMenu | _WSThickFrame | _WSMinimizeBox | _WSMaximizeBox)
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

type _MinMaxInfo struct {
	ptReserved     _Point
	ptMaxSize      _Point
	ptMaxPosition  _Point
	ptMinTrackSize _Point
	ptMaxTrackSize _Point
}

func init() {
	runtime.LockOSThread()

	r, _, _ := ole32CoInitializeEx.Call(0, 2)
	if int(r) < 0 {
		log.Printf("Warning: CoInitializeEx call failed: E=%08x", r)
	}
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
	maxsz      _Point
	minsz      _Point
	m          sync.Mutex
	bindings   map[string]interface{}
	dispatchq  []func()
}

func newchromiumedge(msgcb func(string)) *chromiumedge {
	e := &chromiumedge{msgcb: msgcb}
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
	res, err := createCoreWebView2EnvironmentWithOptions(nil, windows.StringToUTF16Ptr(dataPath), 0, e.envCompleted)
	if err != nil {
		log.Printf("Error calling Webview2Loader: %v", err)
		return false
	} else if res != 0 {
		log.Printf("Result: %08x", res)
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
	if int64(res) < 0 {
		log.Fatalf("Creating environment failed with %08x", res)
	}
	env.vtbl.CreateCoreWebView2Controller.Call(
		uintptr(unsafe.Pointer(env)),
		e.hwnd,
		uintptr(unsafe.Pointer(e.controllerCompleted)),
	)
	return 0
}

func (e *chromiumedge) ControllerCompleted(res uintptr, controller *iCoreWebView2Controller) uintptr {
	if int64(res) < 0 {
		log.Fatalf("Creating controller failed with %08x", res)
	}
	controller.vtbl.AddRef.Call(uintptr(unsafe.Pointer(controller)))
	e.controller = controller

	var token _EventRegistrationToken
	controller.vtbl.GetCoreWebView2.Call(
		uintptr(unsafe.Pointer(controller)),
		uintptr(unsafe.Pointer(&e.webview)),
	)
	e.webview.vtbl.AddRef.Call(
		uintptr(unsafe.Pointer(e.webview)),
	)
	e.webview.vtbl.AddWebMessageReceived.Call(
		uintptr(unsafe.Pointer(e.webview)),
		uintptr(unsafe.Pointer(e.webMessageReceived)),
		uintptr(unsafe.Pointer(&token)),
	)
	e.webview.vtbl.AddPermissionRequested.Call(
		uintptr(unsafe.Pointer(e.webview)),
		uintptr(unsafe.Pointer(e.permissionRequested)),
		uintptr(unsafe.Pointer(&token)),
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
	w.bindings = map[string]interface{}{}
	w.browser = newchromiumedge(w.msgcb)
	w.mainthread, _, _ = kernel32GetCurrentThreadID.Call()
	if !w.Create(debug, window) {
		return nil
	}
	return w
}

type rpcMessage struct {
	ID     int               `json:"id"`
	Method string            `json:"method"`
	Params []json.RawMessage `json:"params"`
}

func jsString(v interface{}) string { b, _ := json.Marshal(v); return string(b) }

func (w *webview) msgcb(msg string) {
	d := rpcMessage{}
	if err := json.Unmarshal([]byte(msg), &d); err != nil {
		log.Printf("invalid RPC message: %v", err)
		return
	}

	id := strconv.Itoa(d.ID)
	if res, err := w.callbinding(d); err != nil {
		w.Dispatch(func() {
			w.Eval("window._rpc[" + id + "].reject(" + jsString(err.Error()) + "); window._rpc[" + id + "] = undefined")
		})
	} else if b, err := json.Marshal(res); err != nil {
		w.Dispatch(func() {
			w.Eval("window._rpc[" + id + "].reject(" + jsString(err.Error()) + "); window._rpc[" + id + "] = undefined")
		})
	} else {
		w.Dispatch(func() {
			w.Eval("window._rpc[" + id + "].resolve(" + string(b) + "); window._rpc[" + id + "] = undefined")
		})
	}
}

func (w *webview) callbinding(d rpcMessage) (interface{}, error) {
	w.m.Lock()
	f, ok := w.bindings[d.Method]
	w.m.Unlock()
	if !ok {
		return nil, nil
	}

	v := reflect.ValueOf(f)
	isVariadic := v.Type().IsVariadic()
	numIn := v.Type().NumIn()
	if (isVariadic && len(d.Params) < numIn-1) || (!isVariadic && len(d.Params) != numIn) {
		return nil, errors.New("function arguments mismatch")
	}
	args := []reflect.Value{}
	for i := range d.Params {
		var arg reflect.Value
		if isVariadic && i >= numIn-1 {
			arg = reflect.New(v.Type().In(numIn - 1).Elem())
		} else {
			arg = reflect.New(v.Type().In(i))
		}
		if err := json.Unmarshal(d.Params[i], arg.Interface()); err != nil {
			return nil, err
		}
		args = append(args, arg.Elem())
	}

	errorType := reflect.TypeOf((*error)(nil)).Elem()
	res := v.Call(args)
	switch len(res) {
	case 0:
		// No results from the function, just return nil
		return nil, nil

	case 1:
		// One result may be a value, or an error
		if res[0].Type().Implements(errorType) {
			if res[0].Interface() != nil {
				return nil, res[0].Interface().(error)
			}
			return nil, nil
		}
		return res[0].Interface(), nil

	case 2:
		// Two results: first one is value, second is error
		if !res[1].Type().Implements(errorType) {
			return nil, errors.New("second return value must be an error")
		}
		if res[1].Interface() == nil {
			return res[0].Interface(), nil
		}
		return res[0].Interface(), res[1].Interface().(error)

	default:
		return nil, errors.New("unexpected number of return values")
	}
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
			lpmmi := (*_MinMaxInfo)(unsafe.Pointer(lp))
			if w.maxsz.x > 0 && w.maxsz.y > 0 {
				lpmmi.ptMaxSize = w.maxsz
				lpmmi.ptMaxTrackSize = w.maxsz
			}
			if w.minsz.x > 0 && w.minsz.y > 0 {
				lpmmi.ptMinTrackSize = w.minsz
			}
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
			w.m.Lock()
			q := append([]func(){}, w.dispatchq...)
			w.dispatchq = []func(){}
			w.m.Unlock()
			for _, v := range q {
				v()
			}
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

func (w *webview) SetSize(width int, height int, hints Hint) {
	index := _GWLStyle
	style, _, _ := user32GetWindowLongPtrW.Call(w.hwnd, uintptr(index))
	if hints == HintFixed {
		style &^= (_WSThickFrame | _WSMaximizeBox)
	} else {
		style |= (_WSThickFrame | _WSMaximizeBox)
	}
	user32SetWindowLongPtrW.Call(w.hwnd, uintptr(index), style)

	if hints == HintMax {
		w.maxsz.x = int32(width)
		w.maxsz.y = int32(height)
	} else if hints == HintMin {
		w.minsz.x = int32(width)
		w.minsz.y = int32(height)
	} else {
		r := _Rect{}
		r.Left = 0
		r.Top = 0
		r.Right = int32(width)
		r.Bottom = int32(height)
		user32AdjustWindowRect.Call(uintptr(unsafe.Pointer(&r)), _WSOverlappedWindow, 0)
		user32SetWindowPos.Call(
			w.hwnd, 0, uintptr(r.Left), uintptr(r.Top), uintptr(r.Right-r.Left), uintptr(r.Bottom-r.Top),
			_SWPNoZOrder|_SWPNoActivate|_SWPNoMove|_SWPFrameChanged)
		w.browser.Resize()
	}
}

func (w *webview) Init(js string) {
	w.browser.Init(js)
}

func (w *webview) Eval(js string) {
	w.browser.Eval(js)
}

func (w *webview) Dispatch(f func()) {
	w.m.Lock()
	w.dispatchq = append(w.dispatchq, f)
	w.m.Unlock()
	user32PostThreadMessageW.Call(w.mainthread, _WMApp, 0, 0)
}

func (w *webview) Bind(name string, f interface{}) error {
	v := reflect.ValueOf(f)
	if v.Kind() != reflect.Func {
		return errors.New("only functions can be bound")
	}
	if n := v.Type().NumOut(); n > 2 {
		return errors.New("function may only return a value or a value+error")
	}
	w.m.Lock()
	w.bindings[name] = f
	w.m.Unlock()

	w.Init("(function() { var name = " + jsString(name) + ";" + `
		var RPC = window._rpc = (window._rpc || {nextSeq: 1});
		window[name] = function() {
		  var seq = RPC.nextSeq++;
		  var promise = new Promise(function(resolve, reject) {
			RPC[seq] = {
			  resolve: resolve,
			  reject: reject,
			};
		  });
		  window.external.invoke(JSON.stringify({
			id: seq,
			method: name,
			params: Array.prototype.slice.call(arguments),
		  }));
		  return promise;
		}
	})()`)

	return nil
}
