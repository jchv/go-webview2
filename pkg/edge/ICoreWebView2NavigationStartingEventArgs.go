package edge

type _ICoreWebView2NavigationStartingEventArgsVtbl struct {
	_IUnknownVtbl
	GetUri                             ComProc
	GetIsUserInitiated                 ComProc
	GetIsRedirected                    ComProc
	GetRequestHeaders                  ComProc
	GetAdditionalAllowedFrameAncestors ComProc
	PutCancel                          ComProc
	GetNavigationId                    ComProc
}

type ICoreWebView2NavigationStartingEventArgs struct {
	vtbl *_ICoreWebView2NavigationStartingEventArgsVtbl
}

func (i *ICoreWebView2NavigationStartingEventArgs) AddRef() uintptr {
	r, _, _ := i.vtbl.AddRef.Call()
	return r
}
