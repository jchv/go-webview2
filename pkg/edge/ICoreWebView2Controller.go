package edge

import (
	"github.com/jchv/go-webview2/internal/w32"
	"golang.org/x/sys/windows"
	"unsafe"
)

type _ICoreWebView2ControllerVtbl struct {
	_IUnknownVtbl
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
	vtbl *_ICoreWebView2ControllerVtbl
}

func (i *iCoreWebView2Controller) AddRef() uintptr {
	return i.AddRef()
}

func (i *iCoreWebView2Controller) GetBounds() (*w32.Rect, error) {
	var err error
	var bounds *w32.Rect
	_, _, err = i.vtbl.GetBounds.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&bounds)),
	)
	if err != windows.ERROR_SUCCESS {
		return nil, err
	}
	return bounds, nil
}

func (i *iCoreWebView2Controller) PutBounds(bounds w32.Rect) error {
	var err error

	_, _, err = i.vtbl.PutBounds.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&bounds)),
	)
	if err != windows.ERROR_SUCCESS {
		return err
	}
	return nil
}

func (i *iCoreWebView2Controller) AddAcceleratorKeyPressed(eventHandler *ICoreWebView2AcceleratorKeyPressedEventHandler, token *_EventRegistrationToken) error {
	var err error
	_, _, err = i.vtbl.AddAcceleratorKeyPressed.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(eventHandler)),
		uintptr(unsafe.Pointer(&token)),
	)
	if err != windows.ERROR_SUCCESS {
		return err
	}
	return nil
}

/*

func (i *iCoreWebView2Controller) GetIsVisible() (bool, error) {
	var err error
	var isVisible bool
	_, _, err = i.vtbl.GetIsVisible.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&isVisible)),
	)
	if err != windows.ERROR_SUCCESS {
		return false, err
	}
	return isVisible, nil
}

func (i *iCoreWebView2Controller) PutIsVisible(isVisible bool) error {
	var err error

	_, _, err = i.vtbl.PutIsVisible.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&isVisible)),
	)
	if err != windows.ERROR_SUCCESS {
		return err
	}
	return nil
}
func (i *iCoreWebView2Controller) GetZoomFactor() (float64, error) {
	var err error
	var zoomFactor float64
	_, _, err = i.vtbl.GetZoomFactor.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&zoomFactor)),
	)
	if err != windows.ERROR_SUCCESS {
		return 0.0, err
	}
	return zoomFactor, nil
}

func (i *iCoreWebView2Controller) PutZoomFactor(zoomFactor float64) error {
	var err error

	_, _, err = i.vtbl.PutZoomFactor.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&zoomFactor)),
	)
	if err != windows.ERROR_SUCCESS {
		return err
	}
	return nil
}

func (i *iCoreWebView2Controller) AddZoomFactorChanged(eventHandler *ICoreWebView2ZoomFactorChangedEventHandler) (*EventRegistrationToken, error) {
	var err error
	var token *EventRegistrationToken
	_, _, err = i.vtbl.AddZoomFactorChanged.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(eventHandler)),
		uintptr(unsafe.Pointer(&token)),
	)
	if err != windows.ERROR_SUCCESS {
		return nil, err
	}
	return token, nil
}

func (i *iCoreWebView2Controller) RemoveZoomFactorChanged(token EventRegistrationToken) error {
	var err error

	_, _, err = i.vtbl.RemoveZoomFactorChanged.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&token)),
	)
	if err != windows.ERROR_SUCCESS {
		return err
	}
	return nil
}

func (i *iCoreWebView2Controller) SetBoundsAndZoomFactor(bounds RECT, zoomFactor float64) error {
	var err error

	_, _, err = i.vtbl.SetBoundsAndZoomFactor.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&bounds)),
		uintptr(unsafe.Pointer(&zoomFactor)),
	)
	if err != windows.ERROR_SUCCESS {
		return err
	}
	return nil
}

func (i *iCoreWebView2Controller) MoveFocus(reason COREWEBVIEW2_MOVE_FOCUS_REASON) error {
	var err error

	_, _, err = i.vtbl.MoveFocus.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&reason)),
	)
	if err != windows.ERROR_SUCCESS {
		return err
	}
	return nil
}

func (i *iCoreWebView2Controller) AddMoveFocusRequested(eventHandler *ICoreWebView2MoveFocusRequestedEventHandler) (*EventRegistrationToken, error) {
	var err error
	var token *EventRegistrationToken
	_, _, err = i.vtbl.AddMoveFocusRequested.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(eventHandler)),
		uintptr(unsafe.Pointer(&token)),
	)
	if err != windows.ERROR_SUCCESS {
		return nil, err
	}
	return token, nil
}

func (i *iCoreWebView2Controller) RemoveMoveFocusRequested(token EventRegistrationToken) error {
	var err error

	_, _, err = i.vtbl.RemoveMoveFocusRequested.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&token)),
	)
	if err != windows.ERROR_SUCCESS {
		return err
	}
	return nil
}

func (i *iCoreWebView2Controller) AddGotFocus(eventHandler *ICoreWebView2FocusChangedEventHandler) (*EventRegistrationToken, error) {
	var err error
	var token *EventRegistrationToken
	_, _, err = i.vtbl.AddGotFocus.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(eventHandler)),
		uintptr(unsafe.Pointer(&token)),
	)
	if err != windows.ERROR_SUCCESS {
		return nil, err
	}
	return token, nil
}

func (i *iCoreWebView2Controller) RemoveGotFocus(token EventRegistrationToken) error {
	var err error

	_, _, err = i.vtbl.RemoveGotFocus.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&token)),
	)
	if err != windows.ERROR_SUCCESS {
		return err
	}
	return nil
}

func (i *iCoreWebView2Controller) AddLostFocus(eventHandler *ICoreWebView2FocusChangedEventHandler) (*EventRegistrationToken, error) {
	var err error
	var token *EventRegistrationToken
	_, _, err = i.vtbl.AddLostFocus.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(eventHandler)),
		uintptr(unsafe.Pointer(&token)),
	)
	if err != windows.ERROR_SUCCESS {
		return nil, err
	}
	return token, nil
}

func (i *iCoreWebView2Controller) RemoveLostFocus(token EventRegistrationToken) error {
	var err error

	_, _, err = i.vtbl.RemoveLostFocus.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&token)),
	)
	if err != windows.ERROR_SUCCESS {
		return err
	}
	return nil
}
*/

/*
func (i *iCoreWebView2Controller) RemoveAcceleratorKeyPressed(token EventRegistrationToken) error {
	var err error

	_, _, err = i.vtbl.RemoveAcceleratorKeyPressed.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&token)),
	)
	if err != windows.ERROR_SUCCESS {
		return err
	}
	return nil
}

func (i *iCoreWebView2Controller) GetParentWindow() (*HWND, error) {
	var err error
	var parentWindow *HWND
	_, _, err = i.vtbl.GetParentWindow.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&parentWindow)),
	)
	if err != windows.ERROR_SUCCESS {
		return nil, err
	}
	return parentWindow, nil
}

func (i *iCoreWebView2Controller) PutParentWindow(parentWindow HWND) error {
	var err error

	_, _, err = i.vtbl.PutParentWindow.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&parentWindow)),
	)
	if err != windows.ERROR_SUCCESS {
		return err
	}
	return nil
}

func (i *iCoreWebView2Controller) NotifyParentWindowPositionChanged() error {
	var err error

	_, _, err = i.vtbl.NotifyParentWindowPositionChanged.Call(
		uintptr(unsafe.Pointer(i)),
	)
	if err != windows.ERROR_SUCCESS {
		return err
	}
	return nil
}

func (i *iCoreWebView2Controller) Close() error {
	var err error

	_, _, err = i.vtbl.Close.Call(
		uintptr(unsafe.Pointer(i)),
	)
	if err != windows.ERROR_SUCCESS {
		return err
	}
	return nil
}

func (i *iCoreWebView2Controller) GetCoreWebView2() (*ICoreWebView2, error) {
	var err error
	var coreWebView2 *ICoreWebView2
	_, _, err = i.vtbl.GetCoreWebView2.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&coreWebView2)),
	)
	if err != windows.ERROR_SUCCESS {
		return nil, err
	}
	return coreWebView2, nil
}
*/
