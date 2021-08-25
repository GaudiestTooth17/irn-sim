package sim

import (
	"math/rand"

	"gonum.org/v1/gonum/mat"
)

func SimForSurvivalRate(M *mat.Dense,
	sir0 SIR,
	disease Disease,
	behavior Behavior,
	maxSteps int,
	rng *rand.Rand) float64 {

	sirs := Simulate(M, sir0, disease, behavior, maxSteps, rng)
	lastSIR := sirs[len(sirs)-1]
	return float64(lastSIR.NumSusceptbile()) / float64(len(lastSIR.S))
}
