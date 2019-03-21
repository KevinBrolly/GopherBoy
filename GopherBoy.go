package main

import (
	"GopherBoy/Gameboy"
	"runtime"
)

func main() {
	runtime.LockOSThread()

	gameboy := Gameboy.NewGameboy()
	gameboy.Run()
}
