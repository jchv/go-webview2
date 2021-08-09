package edge

type _ICoreWebView2WebResourceResponseVtbl struct {
	_IUnknownVtbl
	GetContent      ComProc
	PutContent      ComProc
	GetHeaders      ComProc
	GetStatusCode   ComProc
	PutStatusCode   ComProc
	GetReasonPhrase ComProc
	PutReasonPhrase ComProc
}

type ICoreWebView2WebResourceResponse struct {
	vtbl *_ICoreWebView2WebResourceResponseVtbl
}

func (i *ICoreWebView2WebResourceResponse) AddRef() uintptr {
	return i.AddRef()
}

/*
func (i *ICoreWebView2WebResourceResponse) GetContent() (*IStream, error) {
	var err error
	var content *IStream
	_, _, err = i.vtbl.GetContent.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&content)),
	)
	if err != windows.ERROR_SUCCESS {
		return nil, err
	}
	return content, nil
}

func (i *ICoreWebView2WebResourceResponse) PutContent(content *IStream) error {
	var err error

	_, _, err = i.vtbl.PutContent.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(content)),
	)
	if err != windows.ERROR_SUCCESS {
		return err
	}
	return nil
}

func (i *ICoreWebView2WebResourceResponse) GetHeaders() (*ICoreWebView2HttpResponseHeaders, error) {
	var err error
	var headers *ICoreWebView2HttpResponseHeaders
	_, _, err = i.vtbl.GetHeaders.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&headers)),
	)
	if err != windows.ERROR_SUCCESS {
		return nil, err
	}
	return headers, nil
}

func (i *ICoreWebView2WebResourceResponse) GetStatusCode() (int, error) {
	var err error
	var statusCode int
	_, _, err = i.vtbl.GetStatusCode.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(statusCode),
	)
	if err != windows.ERROR_SUCCESS {
		return 0, err
	}
	return statusCode, nil
}

func (i *ICoreWebView2WebResourceResponse) PutStatusCode(statusCode int) error {
	var err error

	_, _, err = i.vtbl.PutStatusCode.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(statusCode),
	)
	if err != windows.ERROR_SUCCESS {
		return err
	}
	return nil
}

func (i *ICoreWebView2WebResourceResponse) GetReasonPhrase() (string, error) {
	var err error
	// Create *uint16 to hold result
	var _reasonPhrase *uint16
	_, _, err = i.vtbl.GetReasonPhrase.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_reasonPhrase)),
	)
	if err != windows.ERROR_SUCCESS {
		return "", err
	} // Get result and cleanup
	reasonPhrase := windows.UTF16PtrToString(_reasonPhrase)
	windows.CoTaskMemFree(unsafe.Pointer(_reasonPhrase))
	return reasonPhrase, nil
}

func (i *ICoreWebView2WebResourceResponse) PutReasonPhrase(reasonPhrase string) error {
	var err error
	// Convert string 'reasonPhrase' to *uint16
	_reasonPhrase, err := windows.UTF16PtrFromString(reasonPhrase)
	if err != nil {
		return err
	}

	_, _, err = i.vtbl.PutReasonPhrase.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_reasonPhrase)),
	)
	if err != windows.ERROR_SUCCESS {
		return err
	}
	return nil
}
*/
