package sim

import (
	"fmt"

	"math/rand"

	"github.com/GaudiestTooth17/irn-sim/network"
	"github.com/GaudiestTooth17/irn-sim/sets"
	"gonum.org/v1/gonum/mat"
)

type Behavior interface {
	Name() string
	UpdateConnections(D *mat.Dense, M *mat.Dense, timeStep int, sir SIR) *mat.Dense
}

type SimplePressureBehavior struct {
	radius             int
	net                *network.AdjacencyList
	pressure           []float64
	flickerProbability float64
	rng                *rand.Rand
}

func NewSimplePressureBehavior(net *network.AdjacencyList,
	rng *rand.Rand,
	radius int,
	flickerProbability float64) SimplePressureBehavior {

	return SimplePressureBehavior{
		radius:             radius,
		net:                net,
		pressure:           make([]float64, net.N()),
		flickerProbability: flickerProbability,
		rng:                rng,
	}
}

func (b SimplePressureBehavior) Name() string {
	return fmt.Sprintf("SimplePressure(radius=%d, flicker_probability=%f)",
		b.radius, b.flickerProbability)
}

func (b SimplePressureBehavior) UpdateConnections(D *mat.Dense, M *mat.Dense, timeStep int, sir SIR) *mat.Dense {
	infectiousAgents := sir.InfectiousAgents()
	if len(infectiousAgents) > 0 {
		// populate pressuredAgents by finding all the agents in pressure range
		// of the infectious agents
		pressuredAgents := sets.EmptyIntSet()
		for agent := range infectiousAgents {
			agentsInRadius := b.net.NodesWithin(int64(agent), b.radius)
			pressuredAgents = pressuredAgents.Union(agentsInRadius)
		}
		// apply pressure
		for agent := range pressuredAgents {
			b.pressure[agent] += 1
		}
	}

	recoveredAgents := sir.RecoveredAgents()
	if len(recoveredAgents) > 0 {
		// populate unpressuredAgents by finding all the agents in pressure range
		// of the recovered agents
		unpressuredAgents := sets.EmptyIntSet()
		for agent := range recoveredAgents {
			agentsInRange := b.net.NodesWithin(int64(agent), b.radius)
			unpressuredAgents = unpressuredAgents.Union(agentsInRange)
		}
		// relieve pressure
		for agent := range unpressuredAgents {
			b.pressure[agent] -= 1
		}
	}

	// for all the agents currently experiencing pressure, get a random value
	// to determine if they will flicker
	flickeringAgents := sets.EmptyIntSet()
	for agent, pressureValue := range b.pressure {
		if pressureValue > 0 && b.rng.Float64() < b.flickerProbability {
			flickeringAgents.Add(agent)
		}
	}

	R := mat.DenseCopyOf(M)
	// turn off the edges connected to a flickering agent
	R.Apply(func(i, j int, v float64) float64 {
		if flickeringAgents.Contains(i) || flickeringAgents.Contains(j) {
			return 0
		}
		return v
	}, R)
	return R
}
