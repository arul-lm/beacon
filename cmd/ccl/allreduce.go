package main

import (
	"fmt"
	"github.com/arul-lm/beacon/ccl"
	tsim "github.com/arul-lm/triosim/networkmodel"
	"gitlab.com/akita/akita/v3/sim"
	"math/rand"
)

type trafficMsg struct {
	sim.MsgMeta
}

func (m *trafficMsg) Meta() *sim.MsgMeta {
	return &m.MsgMeta
}

func GenerateMsgs(xpus []*ccl.Xpu, n uint64) []sim.Msg {
	msgs := make([]sim.Msg, n)
	for i := uint64(0); i < n; i++ {
		srcAgentID := rand.Intn(len(xpus))
		srcAgent := xpus[srcAgentID]
		srcPortID := rand.Intn(len(srcAgent.Ports()))
		srcPort := srcAgent.Ports()[srcPortID]

		dstAgentID := rand.Intn(len(xpus))
		for dstAgentID == srcAgentID {
			dstAgentID = rand.Intn(len(xpus))
		}

		dstAgent := xpus[dstAgentID]
		dstPortID := rand.Intn(len(dstAgent.Ports()))
		dstPort := dstAgent.Ports()[dstPortID]
		msg := &trafficMsg{}
		msg.Meta().ID = sim.GetIDGenerator().Generate()
		msg.Src = srcPort
		msg.Dst = dstPort
		// msg.TrafficBytes = rand.Intn(4096)
		msg.TrafficBytes = 64 * 1e9
		srcAgent.AddMsg(msg)
		msgs[i] = msg
	}
	for _, x := range xpus {
		fmt.Println("x:", x.MsgsInFlight())
	}
	return msgs
}

func main() {
	engine := sim.NewSerialEngine()
	freq := sim.Freq(1)
	numXpus := 2
	ontw := tsim.NewOpticalNetworkModel(engine, engine, 1, 20)
	xpus := make([]*ccl.Xpu, numXpus)
	for x := range numXpus {
		xpu := ccl.NewXpu(engine, freq, fmt.Sprintf("GPU%d", x), 1)
		ontw.PlugIn(xpu.Ports()[0], 1)
		xpu.TickLater(0)
		xpus[x] = xpu
	}
	ontw.InitHardwareNetwork("mesh", 1, numXpus)
	ontw.InitLogicalNetwork("ring", 1, numXpus)
	GenerateMsgs(xpus, 100)
	err := engine.Run()
	if err != nil {
		panic(err)
	}
	fmt.Println("Sim Time:", engine.CurrentTime())
	for _, x := range(xpus) {
		x.ReportBW(engine.CurrentTime())
	}
}
