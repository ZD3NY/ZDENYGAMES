//go:build js

package main

import "syscall/js"

func exitGame() {
	js.Global().Get("location").Set("href", "/")
}
