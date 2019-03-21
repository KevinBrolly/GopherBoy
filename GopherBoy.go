package main

import (
	"GopherBoy/Gameboy"
	"os"
	"runtime"
)

func main() {
	runtime.LockOSThread()

	gameboy := Gameboy.NewGameboy()

	rom := os.Args[1]
	gameboy.Cartridge.LoadCartridgeData(rom)

	gameboy.Run()
}
