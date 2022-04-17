package main

import (
	"log"

	"github.com/CanPacis/go-webview2"
)

func main() {
	w := webview2.NewWithOptions(webview2.WebViewOptions{
		Debug:     true,
		AutoFocus: true,
		WindowOptions: webview2.WindowOptions{
			Title: "Minimal webview example",
		},
	})
	if w == nil {
		log.Fatalln("Failed to load webview.")
	}
	defer w.Destroy()
	w.SetSize(800, 600, webview2.HintFixed)
	w.Navigate("https://en.m.wikipedia.org/wiki/Main_Page")
	w.Run()
}
