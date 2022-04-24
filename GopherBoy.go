package main

import (
	"image"
	"image/draw"
	"os"

	"github.com/kevinbrolly/GopherBoy/gameboy"

	"github.com/veandco/go-sdl2/sdl"
)

func main() {
	window := NewSDL2Window("Gameboy", 640, 576)

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
		int32(w.Width), int32(w.Height), sdl.WINDOW_ALLOW_HIGHDPI)
	if err != nil {
		panic(err)
	}

	_, err = sdl.CreateRenderer(w.window, -1, 0)
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
	renderer, err := w.window.GetRenderer()
	if err != nil {
		panic(err)
	}

	surface, err := sdl.CreateRGBSurface(0, 160, 144, 32, 0, 0, 0, 0)
	draw.Draw(surface, surface.Bounds(), buffer, image.Point{}, draw.Src)

	texture,err := renderer.CreateTextureFromSurface(surface);
	if err != nil {
		panic(err)
	}
	
	renderer.Copy(texture, nil, nil);
	renderer.Present();
}
