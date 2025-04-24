package main

import (
	"fmt"
	tsim "github.com/arul-lm/triosim/networkmodel"
	"github.com/beacon/ccl"
	"gitlab.com/akita/akita/v3/sim"
)

func main() {
	// tsim.NewOpticalNetworkModel(${1:es sim.EventScheduler}, ${2:tt sim.TimeTeller}, ${3:maxNumWaveGuidesPerNode int}, ${4:establishLatency sim.VTimeInSec})
	engine := sim.NewSerialEngine()
	xpu1 := ccl.NewXpu(engine, engine, "xpu1", 7)
	fmt.Println("Beacon - Allreduce", xpu1)
}
