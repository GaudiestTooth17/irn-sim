package sim

import (
	"math/rand"

	"github.com/GaudiestTooth17/irn-sim/network"
	"gonum.org/v1/gonum/mat"
)

func MultiSimForSurvivalRate(M *mat.Dense,
	sir0 SIR,
	disease Disease,
	behavior Behavior,
	maxSteps int,
	rng *rand.Rand,
	numSims int) []float64 {

	survivalRates := make([]float64, numSims)
	for i := range survivalRates {
		survivalRates[i] = SimForSurvivalRate(M, sir0, disease, behavior, maxSteps, rng)
	}
	return survivalRates
}

// Call MultiSimForSurvivalRate and send the return value through the channel.
func msfsrReturnToChan(M *mat.Dense,
	sir0 SIR,
	disease Disease,
	behavior Behavior,
	maxSteps int,
	rng *rand.Rand,
	numSims int,
	resultChan chan<- []float64) {

	survivalRates := MultiSimForSurvivalRate(M, sir0, disease, behavior, maxSteps, rng, numSims)
	resultChan <- survivalRates
}

// Calls MultiSimForSurvivalRate on each matrix in Ms. A new *rand.Rand is constructed
// with the provided seed each time MultiSimForSurvivalRate is called. The returned slice
// contains the survival rates in arbitrary order, thus this function should be used
// for classes of networks where the goal is to generate a distribution of outcomes
// for all the classes as a whole.
func SimOnManyNetworksForSurvivalRate(nets []*network.AdjacencyList,
	makeSir0 func(N int, numToInfect int, rng *rand.Rand) SIR,
	disease Disease,
	makeBehavior func(*network.AdjacencyList, *rand.Rand) Behavior,
	maxSteps int,
	seed int64,
	numSimsPerNet int) []float64 {

	// start goroutines
	// job chan
	jobs := make(chan int, 4)
	go func() {
		for jobID := range nets {
			jobs <- jobID

		}
		close(jobs)
	}()
	// read values from the job channel to control how many goroutines run at once
	survivalRateChan := make(chan []float64, 4)
	defer close(survivalRateChan)
	for jobID := range jobs {
		net := nets[jobID]
		M := net.M()
		rng := rand.New(rand.NewSource(seed))
		sir0 := makeSir0(net.N(), 1, rng)
		behavior := makeBehavior(net, rng)
		go msfsrReturnToChan(M, sir0, disease, behavior, maxSteps, rng,
			numSimsPerNet, survivalRateChan)
	}

	// read data from channel
	survivalRates := make([]float64, numSimsPerNet*len(nets))
	for i := range nets {
		partialSurvivalRates := <-survivalRateChan
		// copy data
		for j, rate := range partialSurvivalRates {
			survivalRates[i*numSimsPerNet+j] = rate
		}
	}

	return survivalRates
}
