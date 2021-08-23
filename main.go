package main

import (
	"fmt"
	"os"

	fio "github.com/GaudiestTooth17/irn-sim/fileio"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <network>\n", os.Args[0])
		return
	}
	networkPath := os.Args[1]

	nets := fio.ReadClass(networkPath)
	fmt.Printf("Read %d networks.\n", len(nets))

	// // set up the parameters
	// rng := rand.New(rand.NewSource(0))
	// net := fio.ReadFile(networkPath)
	// disease := sim.Disease{DaysInfectious: 4, TransProb: .2}
	// behavior := sim.NewSimplePressureBehavior(net, rng, 2, .25)
	// // behavior := sim.StaticBehavior{}
	// sir0 := sim.MakeSir0(net.N(), 1, rng)
	// // run a simulation
	// // result := sim.Simulate(net.M(), sir0, disease, behavior, 300, rng)
	// // fmt.Println(sim.GetSurvivalPercentage(result))
	// fmt.Println(runSimBatch(net.M(), sir0, disease, behavior, 100, rng, 500))
}
