package display

import (
    "github.com/veandco/go-sdl2/sdl"
)

type Window struct {
    Name string

    window *sdl.Window
    surface *sdl.Surface
    texture *sdl.Texture
    renderer *sdl.Renderer

    Width int
    Height int
    Framebuffer []uint32
    QuitFunc func()
}

func NewWindow(name string, width, height int, quitFunc func()) *Window {
    w := &Window{
        Name: name,
        Width: width,
        Height: height,
        Framebuffer: make([]uint32, width*height),
        QuitFunc: quitFunc,
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

    w.renderer, _ = sdl.CreateRenderer(w.window, -1, sdl.RENDERER_SOFTWARE)

    w.texture, _  = w.renderer.CreateTexture(uint32(sdl.PIXELFORMAT_RGBA32), sdl.TEXTUREACCESS_TARGET, int32(w.Width), int32(w.Height))

    return w
}

func (w *Window) Quit() {
    w.texture.Destroy()
    w.renderer.Destroy()
    w.window.Destroy()
    sdl.Quit()
}

func (w *Window) Update() {
    for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
        switch event.(type) {
        case *sdl.QuitEvent:
            w.Quit()
            w.QuitFunc()
        }
    }

    w.texture.UpdateRGBA(nil, w.Framebuffer, w.Width)
    w.renderer.Clear()
    w.renderer.Copy(w.texture, nil, nil)
    w.renderer.Present()
}
