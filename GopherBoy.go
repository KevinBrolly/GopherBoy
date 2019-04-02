package main

import (
	"os"
	"runtime"

	"github.com/kevinbrolly/GopherBoy/gameboy"
)

func main() {
	runtime.LockOSThread()

	gameboy := gameboy.NewGameboy()

	rom := os.Args[1]
	gameboy.Cartridge.LoadCartridgeData(rom)

	gameboy.Run()
}
