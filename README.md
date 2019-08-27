# GopherBoy

GopherBoy is a Nintendo GameBoy (DMG) emulator written in Go.

It was primarily build as a development excerise to learn more about emulaton. Please feel free to contribute if you're interested in GameBoy emulator development.

Currently GopherBoy will run most DMG games with sound, Color support is coming soon.

## Try it out

With go installed, you can build and run GopherBoy as follows:

```sh
git clone https://github.com/KevinBrolly/GopherBoy.git
cd GopherBoy
go run GopherBoy.go "<path_to_rom>"
```

GopherBoy uses [SDL2](https://www.libsdl.org/) for control binding and graphics, you must have SLD2 installed to use GopherBoy.

## Controls
<kbd>&larr;</kbd> <kbd>&uarr;</kbd> <kbd>&darr;</kbd> <kbd>&rarr;</kbd> <kbd>A</kbd> <kbd>S</kbd> <kbd>Enter</kbd> <kbd>Backspace</kbd>
