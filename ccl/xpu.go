package ccl

import (
	"gitlab.com/akita/akita/v3/sim"
)

type Xpu struct {
	*sim.TickingComponent
	ports []sim.Port
	msgs []sim.Msg
}

func NewXpu(engine sim.Engine, freq sim.Freq, name string, numPorts int) *Xpu{
	xpu := &Xpu{}
	xpu.TickingComponent = sim.NewTickingComponent(name, engine, freq, xpu)
	return xpu
}
