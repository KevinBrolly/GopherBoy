package Gameboy_test

import (
    "testing"
    "GoBoy/Gameboy"
)


func TestP1(t *testing.T) {
    cases := []struct {
        Name string
        Mode byte
        Before  byte
        After byte
    }{
        {"Mode 0", Gameboy.MODE0, 0x80, 0x80},
        {"Mode 1", Gameboy.MODE1, 0x80, 0x81},
        {"Mode 2", Gameboy.MODE2, 0x80, 0x82},
        {"Mode 3", Gameboy.MODE3, 0x80, 0x83},
    }
    for _, tt := range cases {
        t.Run(tt.Name, func(t *testing.T) {
            tt.Before = Gameboy.SetBit(tt.Before, tt.pos)
            if tt.before != tt.after {
                t.Errorf("got %#x, want %#x", tt.before, tt.after)
            }
        })
    }
}
