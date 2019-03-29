package main

import (
	"GopherBoy/gameboy"
	"os"
	"runtime"
)

func main() {
	runtime.LockOSThread()

	gameboy := gameboy.NewGameboy()

	rom := os.Args[1]
	gameboy.Cartridge.LoadCartridgeData(rom)

	gameboy.Run()
}
