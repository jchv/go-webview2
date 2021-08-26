package edge

type _ICoreWebView2NavigationCompletedEventArgsVtbl struct {
	_IUnknownVtbl
	GetIsSuccess      ComProc
	GetWebErrorStatus ComProc
	GetNavigationId   ComProc
}

type ICoreWebView2NavigationCompletedEventArgs struct {
	vtbl *_ICoreWebView2NavigationCompletedEventArgsVtbl
}

func (i *ICoreWebView2NavigationCompletedEventArgs) AddRef() uintptr {
	return i.AddRef()
}

//
//func (i *ICoreWebView2NavigationCompletedEventArgs) GetIsSuccess() (bool, error) {
//	var err error
//	var isSuccess bool
//	_, _, err = i.vtbl.GetIsSuccess.Call(
//		uintptr(unsafe.Pointer(i)),
//		uintptr(unsafe.Pointer(&isSuccess)),
//	)
//	if err != windows.ERROR_SUCCESS {
//		return false, err
//	}
//	return isSuccess, nil
//}
//
//func (i *ICoreWebView2NavigationCompletedEventArgs) GetWebErrorStatus() (*COREWEBVIEW2_WEB_ERROR_STATUS, error) {
//	var err error
//	var webErrorStatus *COREWEBVIEW2_WEB_ERROR_STATUS
//	_, _, err = i.vtbl.GetWebErrorStatus.Call(
//		uintptr(unsafe.Pointer(i)),
//		uintptr(unsafe.Pointer(&webErrorStatus)),
//	)
//	if err != windows.ERROR_SUCCESS {
//		return nil, err
//	}
//	return webErrorStatus, nil
//}
//
//func (i *ICoreWebView2NavigationCompletedEventArgs) GetNavigationId() (*uint64, error) {
//	var err error
//	var navigationId *uint64
//	_, _, err = i.vtbl.GetNavigationId.Call(
//		uintptr(unsafe.Pointer(i)),
//		uintptr(unsafe.Pointer(&navigationId)),
//	)
//	if err != windows.ERROR_SUCCESS {
//		return nil, err
//	}
//	return navigationId, nil
//}
