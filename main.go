package main

import (
	"fmt"
	"math/rand"
	"os"

	fio "github.com/GaudiestTooth17/irn-sim/fileio"
	"github.com/GaudiestTooth17/irn-sim/network"
	"github.com/GaudiestTooth17/irn-sim/sim"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <network>\n", os.Args[0])
		return
	}
	networkPath := os.Args[1]

	nets := fio.ReadClass(networkPath)
	// set up the parameters
	disease := sim.Disease{DaysInfectious: 4, TransProb: .2}
	makeBehavior := func(net *network.AdjacencyList, rng *rand.Rand) sim.Behavior {
		return sim.NewSimplePressureBehavior(net, rng, 2, .25)
	}
	// behavior := sim.StaticBehavior{}
	makeSIR0 := func(N int, numToInfect int, rng *rand.Rand) sim.SIR {
		return sim.MakeSir0(N, 1, rng)
	}
	// run a simulation
	// result := sim.Simulate(net.M(), sir0, disease, behavior, 300, rng)
	// fmt.Println(sim.GetSurvivalPercentage(result))
	survivalRates := sim.SimOnManyNetworksForSurvivalRate(nets, makeSIR0, disease,
		makeBehavior, 150, 0, 100)
	fmt.Println(len(survivalRates))
}
