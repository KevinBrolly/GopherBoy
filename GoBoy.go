package main

import (
    //"fmt"
    "os"
    "log"
    "flag"
    "runtime/pprof"
    "GoBoy/Gameboy"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")

func main() {
    flag.Parse()
    if *cpuprofile != "" {
        f, err := os.Create(*cpuprofile)
        if err != nil {
            log.Fatal(err)
        }
        pprof.StartCPUProfile(f)
        defer pprof.StopCPUProfile()
    }
    gameboy := Gameboy.NewGameboy()
    //gameboy.LoadROM("cpu_instrs.gb")
    //gameboy.LoadROM("gb-test-roms-master/cpu_instrs/individual/01-special.gb")
    //gameboy.LoadROM("gb-test-roms-master/cpu_instrs/individual/02-interrupts.gb")
    //gameboy.LoadROM("gb-test-roms-master/cpu_instrs/individual/03-op sp,hl.gb")
    //gameboy.LoadROM("gb-test-roms-master/cpu_instrs/individual/04-op r,imm.gb")
    //gameboy.LoadROM("gb-test-roms-master/cpu_instrs/individual/05-op rp.gb")
    //gameboy.LoadROM("gb-test-roms-master/cpu_instrs/individual/06-ld r,r.gb")
    gameboy.LoadROM("gb-test-roms-master/cpu_instrs/individual/07-jr,jp,call,ret,rst.gb")
    //gameboy.LoadROM("gb-test-roms-master/cpu_instrs/individual/08-misc instrs.gb")
    //gameboy.LoadROM("gb-test-roms-master/cpu_instrs/individual/09-op r,r.gb")
    //gameboy.LoadROM("gb-test-roms-master/cpu_instrs/individual/10-bit ops.gb")
    //gameboy.LoadROM("gb-test-roms-master/cpu_instrs/individual/11-op a,(hl).gb")
    //fmt.Printf("%v", gameboy.ROM)
    gameboy.Run()
}
