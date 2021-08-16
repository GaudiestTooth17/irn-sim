package main

import (
	"math/rand"

	"github.com/GaudiestTooth17/irn-sim/sim"
	"gonum.org/v1/gonum/mat"
)

// Return the survival rate of each of the simulations
func runSimBatch(M *mat.Dense,
	sir0 sim.SIR,
	disease sim.Disease,
	behavior sim.Behavior,
	maxSimSteps int,
	rng *rand.Rand,
	numSims int) []float64 {

	survivalRates := make([]float64, numSims)
	for i := range survivalRates {
		survivalRates[i] = sim.GetSurvivalPercentage(
			sim.Simulate(M, sir0, disease, behavior, maxSimSteps, rng))
	}
	return survivalRates
}

// func runSimulationAsync(M *mat.Dense,
// 	sir0 sim.SIR,
// 	disease sim.Disease,
// 	maxSteps int,
// 	rng *rand.Rand,
// 	outChan chan float64) {

// }
