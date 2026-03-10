//go:build js

package main

import "syscall/js"

func submitScore(score, lines int) {
	js.Global().Call("submitScore", score, lines)
}
