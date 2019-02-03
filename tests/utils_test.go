package Gameboy_test

import (
    "testing"
    "GoBoy/Gameboy"
)

var setbittests = []struct {
    pos byte
    before  byte
    after byte
}{
    {7, 0x0, 0x80},
    {6, 0x0, 0x40},
    {5, 0x0, 0x20},
    {4, 0x0, 0x10},
    {3, 0x0, 0x08},
    {0, 0x40, 0x41},
    {1, 0x41, 0x43},
}
func TestSetBit(t *testing.T) {
    for _, tt := range setbittests {
        t.Run(string(tt.pos), func(t *testing.T) {
            tt.before = Gameboy.SetBit(tt.before, tt.pos)
            if tt.before != tt.after {
                t.Errorf("got %#x, want %#x", tt.before, tt.after)
            }
        })
    }
}

var isbitsettests = []struct {
    pos byte
    value  byte
    expected_result bool
}{
    {7, 0x80, true},
    {6, 0x40, true},
    {5, 0x20, true},
    {4, 0x10, true},
    {3, 0x0, false},
    {2, 0x80, false},
}
func TestIsBitSet(t *testing.T) {
    for _, tt := range isbitsettests {
        t.Run(string(tt.pos), func(t *testing.T) {
            result := Gameboy.IsBitSet(tt.value, tt.pos)
            if tt.expected_result != result {
                t.Errorf("got %t, want %t", result, tt.expected_result)
            }
        })
    }
}

var clearbittests = []struct {
    pos byte
    before  byte
    after byte
}{
    {7, 0x80, 0x0},
    {6, 0x40, 0x0},
    {5, 0x20, 0x0},
    {4, 0x10, 0x0},
    {3, 0x08, 0x0},
    {0, 0x41, 0x40},
    {1, 0x43, 0x41},
}
func TestClearBit(t *testing.T) {
    for _, tt := range clearbittests {
        t.Run(string(tt.pos), func(t *testing.T) {
            tt.before = Gameboy.ClearBit(tt.before, tt.pos)
            if tt.before != tt.after {
                t.Errorf("got %#x, want %#x", tt.before, tt.after)
            }
        })
    }
}

var joinbytestests = []struct {
    hb byte
    lb  byte
    expected uint16
}{
    {0x80, 0x03, 0x8003},
    {0x01, 0x90, 0x0190},
}
func TestJoinBytes(t *testing.T) {
    for _, tt := range joinbytestests {
        t.Run(string(tt.expected), func(t *testing.T) {
            result := Gameboy.JoinBytes(tt.hb, tt.lb)
            if result != tt.expected {
                t.Errorf("got %#x, want %#x", result, tt.expected)
            }
        })
    }
}

var splitbytetests = []struct {
    split uint16
    hb byte
    lb  byte
}{
    {0x8003, 0x80, 0x03},
    {0x0190, 0x01, 0x90},
}
func TestSplitBytes(t *testing.T) {
    for _, tt := range splitbytetests {
        t.Run(string(tt.split), func(t *testing.T) {
            result_hb, result_lb := Gameboy.SplitBytes(tt.split)
            if result_hb != tt.hb {
                t.Errorf("got %v, want %v", result_hb, tt.hb)
            }
            if result_lb != tt.lb {
                t.Errorf("got %v, want %v", result_lb, tt.lb)
            }
        })
    }
}
