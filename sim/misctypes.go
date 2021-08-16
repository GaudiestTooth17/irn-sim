package sim

import (
	"math/rand"

	"github.com/GaudiestTooth17/irn-sim/sets"
)

type SIR struct {
	S []int
	I []int
	R []int
}

// Make the initial SIR for a simulation with N agent
func MakeSir0(N int, numToInfect int, rng *rand.Rand) SIR {
	s := make([]int, N)
	i := make([]int, N)
	r := make([]int, N)

	// randomly choose the states
	infectious := make(map[int]bool, numToInfect)
	for len(infectious) < numToInfect {
		infectious[rand.Intn(numToInfect)] = true
	}

	// update arrays to reflect what state each agent is in
	for agent := 0; agent < N; agent++ {
		if infectious[agent] {
			i[agent] = 1
		} else {
			s[agent] = 1
		}
	}

	return SIR{s, i, r}
}

func (sir SIR) NumRemoved() int {
	count := 0
	for _, timeInState := range sir.R {
		if timeInState > 0 {
			count++
		}
	}
	return count
}

func (sir SIR) InfectiousLongerThan(time int) []int {
	nodesInfectiousForLongTime := make([]int, 0)
	for node, timeInState := range sir.I {
		if timeInState > time {
			nodesInfectiousForLongTime = append(nodesInfectiousForLongTime, node)
		}
	}
	return nodesInfectiousForLongTime
}

func (sir SIR) RemovedAgents() sets.IntSet {
	theRecovered := sets.EmptyIntSet()
	for agent, timeInState := range sir.R {
		if timeInState > 0 {
			theRecovered.Add(agent)
		}
	}
	return theRecovered
}

func (sir SIR) InfectiousAgents() sets.IntSet {
	theInfectious := sets.EmptyIntSet()
	for agent, timeInState := range sir.I {
		if timeInState > 0 {
			theInfectious.Add(agent)
		}
	}
	return theInfectious
}

func (sir SIR) SusceptibleAgents() sets.IntSet {
	theSusceptible := sets.EmptyIntSet()
	for agent, timeInState := range sir.S {
		if timeInState > 0 {
			theSusceptible.Add(agent)
		}
	}
	return theSusceptible
}

func (sir SIR) SetTimeRecovered(agents []int, time int) {
	for _, agent := range agents {
		sir.R[agent] = time
	}
}

func (sir SIR) SetTimeInfectious(agents []int, time int) {
	for _, agent := range agents {
		sir.I[agent] = time
	}
}

func (sir SIR) SetTimeSusceptible(agents []int, time int) {
	for _, agent := range agents {
		sir.S[agent] = time
	}
}

func (sir SIR) IncrementPositiveTimes() {
	for i := 0; i < len(sir.S); i++ {
		if sir.S[i] > 0 {
			sir.S[i]++
		}
		if sir.I[i] > 0 {
			sir.I[i]++
		}
		if sir.R[i] > 0 {
			sir.R[i]++
		}
	}
}

func (sir SIR) setNegativeTimesTo1() {
	for i := 0; i < len(sir.S); i++ {
		if sir.S[i] < 0 {
			sir.S[i] = 1
		}
		if sir.I[i] < 0 {
			sir.I[i] = 1
		}
		if sir.R[i] < 0 {
			sir.R[i] = 1
		}
	}
}

func (sir SIR) DiseaseGone() bool {
	// for _, timeInState := range sir.I {
	// 	if timeInState > 0 {
	// 		return false
	// 	}
	// }
	// return true
	return len(sir.InfectiousAgents()) > 0
}

func (sir SIR) Copy() SIR {
	newS := make([]int, len(sir.S))
	for i, v := range sir.S {
		newS[i] = v
	}
	newI := make([]int, len(sir.I))
	for i, v := range sir.I {
		newI[i] = v
	}
	newR := make([]int, len(sir.R))
	for i, v := range sir.R {
		newR[i] = v
	}
	return SIR{newS, newI, newR}
}

type Disease struct {
	DaysInfectious int
	TransProb      float64
}
