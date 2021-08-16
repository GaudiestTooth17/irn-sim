package test

import (
	"math/rand"
	"testing"

	fio "github.com/GaudiestTooth17/irn-sim/fileio"
	"github.com/GaudiestTooth17/irn-sim/sim"
)

func TestSimulation(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	net := fio.ReadFile("../networks/cgg-500.txt")
	disease := sim.Disease{DaysInfectious: 4, TransProb: 1}
	behavior := sim.StaticBehavior{}
	sir0 := sim.MakeSir0(net.N(), 1, rng)
	// run a simulation
	result := sim.Simulate(net.M(), sir0, disease, behavior, 300, rng)
	survivalPercentage := sim.GetSurvivalPercentage(result)
	if survivalPercentage != 0.0 {
		t.Errorf("Expected no agents to survive, but %.3f%% did.", survivalPercentage)
	}
}
