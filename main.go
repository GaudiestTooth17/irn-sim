package main

import (
	"fmt"
	"math/rand"
	"os"

	"github.com/GaudiestTooth17/irn-sim/sim"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <network>\n", os.Args[0])
		return
	}
	networkPath := os.Args[1]

	// set up the parameters
	rng := rand.New(rand.NewSource(0))
	net := readFile(networkPath)
	disease := sim.Disease{DaysInfectious: 4, TransProb: .2}
	behavior := sim.NewSimplePressureBehavior(net, rng, 2, .25)
	sir0 := sim.MakeSir0(net.N(), 1, rng)
	// run a simulation
	result := sim.Simulate(net.M(), sir0, disease, behavior, 300, rng)
	fmt.Println(sim.GetSurvivalPercentage(result))
}
