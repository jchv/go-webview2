package main

import (
	"fmt"
	"log"

	"github.com/jchv/go-webview2"
)

func main() {
	w := webview2.NewWithOptions(webview2.WebViewOptions{
		Debug:     true,
		AutoFocus: true,
		WindowOptions: webview2.WindowOptions{
			Title:  "Minimal webview example",
			Width:  800,
			Height: 600,
			IconId: 2, // icon resource id
			Center: true,
		},
	})
	if w == nil {
		log.Fatalln("Failed to load webview.")
	}
	defer w.Destroy()
	w.SetSize(800, 600, webview2.HintFixed)
	_ = w.NavigationStarting(onNavigationStarting)
	_ = w.NavigationCompleted(onNavigationCompleted)
	w.Navigate("https://en.m.wikipedia.org/wiki/Main_Page")
	w.Run()
}

func onNavigationCompleted(httpStatusCode int32, isSuccess bool, navigationId uint64, webErrorStatus int32) {
	fmt.Println("navigation completed:")
	fmt.Println("http status code: ", httpStatusCode)
	fmt.Println("success: ", isSuccess)
	fmt.Println("navigation id: ", navigationId)
	fmt.Println("web error status: ", webErrorStatus)
}

func onNavigationStarting(additionalAllowedFrameAncestors string, isRedirected bool, isUserInitiated bool, navigationId uint64, uri string) bool {
	fmt.Println("navigation starting:")
	fmt.Println("redirected: ", isRedirected)
	fmt.Println("user initiated: ", isUserInitiated)
	fmt.Println("navigation id: ", navigationId)
	fmt.Println("additional allowed frame ancestors: ", additionalAllowedFrameAncestors)
	fmt.Println("uri: ", uri)
	return false
}
