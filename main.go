package main

import (
	"fmt"
	"os"

	"github.com/GaudiestTooth17/irn-sim/sim"
	"gonum.org/v1/gonum/stat/distuv"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <network>\n", os.Args[0])
		return
	}
	networkPath := os.Args[1]

	// set up the parameters
	rng := distuv.UnitUniform
	net := readFile(networkPath)
	disease := sim.Disease{4, .2}
	behavior := sim.NewSimplePressureBehavior(net, rng, 2, .25)
	sir0 := sim.MakeSir0(1, rng)
	// run a simulation
	result := sim.Simulate(net.M(), sir)
}
