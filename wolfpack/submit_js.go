//go:build js

package main

import "syscall/js"

func submitScore(score, waves int) {
	js.Global().Call("submitScore", score, waves)
}

