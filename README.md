# go-webview2
This is a proof of concept for embedding Webview2 into Go without CGo. It is based on [webview/webview](https://github.com/webview/webview) and provides a compatible API sans some unimplemented functionality (notably, bindings are not implemented.)

## Notice
Because this version doesn't currently have an EdgeHTML fallback, it will not work unless you have a Webview2 runtime installed. In addition, it requires the Webview2Loader DLL in order to function. Adding an EdgeHTML fallback should be technically possible but will likely require much worse hacks since the API is not strictly COM to my knowledge.

## Demo
To run the demo, `git clone` or otherwise download the repo locally, then run `go run ./cmd/demo` from the root. There is a x64 copy of `WebView2Loader.dll` in the root that should get picked up, but you will still need a copy of the [WebView2 runtime](https://developer.microsoft.com/en-us/microsoft-edge/webview2/) installed as well.

If you want to run this on another CPU architecture, you'll need a different copy of WebView2Loader.dll, which can be acquired from the [NuGet package](https://www.nuget.org/packages/Microsoft.Web.WebView2). Note that the bindings have not been tested on x86-32 or AArch64 systems (it may potentially still work; if someone wants to support this send PRs please.)
