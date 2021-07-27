package main

import (
	"fmt"
	"image"
	"image/draw"
	"os"

	"github.com/kevinbrolly/GopherBoy/gameboy"

	"github.com/veandco/go-sdl2/sdl"
)

func main() {
	fmt.Print("Test")
	window := NewSDL2Window("Gameboy", 160, 144)

	gameboy := gameboy.NewGameboy(window)

	rom := os.Args[1]
	gameboy.LoadCartridge(rom)

	gameboy.Run()
}

type SDL2Window struct {
	Name   string
	Width  int
	Height int

	window *sdl.Window
}

func NewSDL2Window(name string, width, height int) *SDL2Window {
	fmt.Print("Ttest")
	w := &SDL2Window{
		Name:   name,
		Width:  width,
		Height: height,
	}

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}

	var err error
	w.window, err = sdl.CreateWindow(w.Name, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		int32(w.Width), int32(w.Height), sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}

	return w
}

func (w *SDL2Window) Quit() {
	w.window.Destroy()
	sdl.Quit()
}

func (w *SDL2Window) DrawFrame(buffer *image.RGBA) {
	surface, err := w.window.GetSurface()
	if err != nil {
		panic(err)
	}

	draw.Draw(surface, surface.Bounds(), buffer, image.Point{}, draw.Src)
	w.window.UpdateSurface()
}
