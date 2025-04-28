package ccl

import (
	"fmt"
	"gitlab.com/akita/akita/v3/sim"
	"log"
)

type Xpu struct {
	*sim.TickingComponent
	ports     []sim.Port
	msgs      []sim.Msg
	sendBytes uint64
	recvBytes uint64
	recvMsgs  map[sim.Msg]bool
}

func (x *Xpu) Ports() []sim.Port {
	return x.ports
}

func (x *Xpu) AddMsg(msg sim.Msg) {
	x.msgs = append(x.msgs, msg)
}
func (x *Xpu) MsgsInFlight() int {
	return len(x.msgs)
}

func (x *Xpu) Tick(now sim.VTimeInSec) bool {
	// name := x.TickingComponent.Name()
	madeProgress := x.send(now)
	// fmt.Println("XPU Tick:", name, len(x.msgs), len(x.recvMsgs), x.sendBytes, x.recvBytes, madeProgress)
	recvRes := x.recv(now)
	madeProgress = recvRes || madeProgress
	// fmt.Println("XPU Tick:", name, len(x.msgs), len(x.recvMsgs), x.sendBytes, x.recvBytes, recvRes, madeProgress)
	return madeProgress
}

func (x *Xpu) send(now sim.VTimeInSec) bool {
	if len(x.msgs) == 0 {
		return false
	}
	msg := x.msgs[0]
	msg.Meta().SendTime = now
	err := msg.Meta().Src.Send(msg)
	if err == nil {
		x.msgs = x.msgs[1:]
		x.sendBytes += uint64(msg.Meta().TrafficBytes)
		return true
	}
	return false
}

func (x *Xpu) recv(now sim.VTimeInSec) bool {
	madeProgress := false
	// name := x.TickingComponent.Name()
	for _, port := range x.ports {
		msg := port.Retrieve(now)
		if msg != nil {
			if msg.Meta().Dst != port {
				panic("Msg delivered at the wrong destination")
			}
			if _, found := x.recvMsgs[msg]; found {
				panic("Msg already received")
			}
			x.recvMsgs[msg] = true
			x.recvBytes += uint64(msg.Meta().TrafficBytes)
			madeProgress = true
			// fmt.Println("Received msg", name, msg)
		}
	}
	return madeProgress
}

func NewXpu(engine sim.Engine, freq sim.Freq, name string, numPorts int) *Xpu {
	x := &Xpu{}
	x.ports = make([]sim.Port, numPorts)
	x.TickingComponent = sim.NewTickingComponent(name, engine, freq, x)
	for i := range numPorts {
		p := sim.NewLimitNumMsgPort(x, 1, fmt.Sprintf("%sPort", name))
		x.ports[i] = p
	}
	x.recvMsgs = make(map[sim.Msg]bool)
	return x
}

func (x *Xpu) ReportBW(now sim.VTimeInSec) {
	log.Printf("xpu %s, send bandwidth %.2f GB/s, recv bandwidth %.2f GB/s",
		x.Name(),
		float64(x.sendBytes)/float64(now)/1e9,
		float64(x.recvBytes)/float64(now)/1e9)
}
