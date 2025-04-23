package edge

type _ICoreWebView2NavigationStartingEventHandlerVtbl struct {
	_IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2NavigationStartingEventHandler struct {
	vtbl *_ICoreWebView2NavigationStartingEventHandlerVtbl
	impl _ICoreWebView2NavigationStartingEventHandlerImpl
}

func _ICoreWebView2NavigationStartingEventHandlerIUnknownQueryInterface(this *ICoreWebView2NavigationStartingEventHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func _ICoreWebView2NavigationStartingEventHandlerIUnknownAddRef(this *ICoreWebView2NavigationStartingEventHandler) uintptr {
	return this.impl.AddRef()
}

func _ICoreWebView2NavigationStartingEventHandlerIUnknownRelease(this *ICoreWebView2NavigationStartingEventHandler) uintptr {
	return this.impl.Release()
}

func _ICoreWebView2NavigationStartingEventHandlerInvoke(this *ICoreWebView2NavigationStartingEventHandler, sender *ICoreWebView2, args *ICoreWebView2NavigationStartingEventArgs) uintptr {
	return this.impl.NavigationStarting(sender, args)
}

type _ICoreWebView2NavigationStartingEventHandlerImpl interface {
	_IUnknownImpl
	NavigationStarting(sender *ICoreWebView2, args *ICoreWebView2NavigationStartingEventArgs) uintptr
}

var _ICoreWebView2NavigationStartingEventHandlerFn = _ICoreWebView2NavigationStartingEventHandlerVtbl{
	_IUnknownVtbl{
		NewComProc(_ICoreWebView2NavigationStartingEventHandlerIUnknownQueryInterface),
		NewComProc(_ICoreWebView2NavigationStartingEventHandlerIUnknownAddRef),
		NewComProc(_ICoreWebView2NavigationStartingEventHandlerIUnknownRelease),
	},
	NewComProc(_ICoreWebView2NavigationStartingEventHandlerInvoke),
}

func newICoreWebView2NavigationStartingEventHandler(impl _ICoreWebView2NavigationStartingEventHandlerImpl) *ICoreWebView2NavigationStartingEventHandler {
	return &ICoreWebView2NavigationStartingEventHandler{
		vtbl: &_ICoreWebView2NavigationStartingEventHandlerFn,
		impl: impl,
	}
}
