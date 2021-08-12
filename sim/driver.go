package sim

import (
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/stat/distuv"
)

func Simulate(M *mat.Dense,
	sir0 SIR,
	disease Disease,
	behavior Behavior,
	maxSteps int,
	rng distuv.Rander) []SIR {

	sirs := make([]SIR, maxSteps)
	sirs[0] = sir0.Copy()
	D := mat.DenseCopyOf(M)
	N, _ := M.Dims()

	for step := 1; step < maxSteps; step++ {
		// get the adjacency matrix to use at this step
		D = behavior.UpdateConnections(D, M, step, sirs[step-1])

		// nextSIR is the workhorse of the simulation because it is responsible
		// for simulating the disease spread
		newSir, statesChanged := nextSIR(sirs[step-1], D, disease, rng)
		sirs[step] = newSir

		// find all the agents that are in the removed state. If that number is N,
		// then the simulation is done.
		if !statesChanged && sirs[step].NumRemoved() == N {
			return sirs[:step]
		}

		// If there aren't any infectious agents, the disease is gone and we
		// can take a short cut to finish the simulation.
		if !statesChanged && sirs[step].DiseaseGone() {
			for i := step; i < maxSteps; i++ {
				sirs[i] = sirs[step].Copy()
			}
			return sirs
		}
	}
	return sirs
}

func nextSIR(oldSIR SIR, M *mat.Dense, disease Disease, rng distuv.Rander) (SIR, bool) {
	sir := oldSIR.Copy()

	// infectious to recovered
	toRFilter := sir.InfectiousLongerThan(disease.DaysInfectious)
	sir.SetTimeRecovered(toRFilter, -1)
	sir.SetTimeInfectious(toRFilter, 0)

	// susceptible to infectious
	iFilter := sir.InfectiousAgents()
	toIProbs := calculateToIProbs(M, disease, iFilter.Values())
	toIFilter := makeToIFilter(sir, toIProbs, rng)
	sir.SetTimeInfectious(toIFilter, -1)
	sir.SetTimeSusceptible(toIFilter, 0)

	sir.IncrementPositiveTimes()
	sir.setNegativeTimesTo1()

	return sir, len(toRFilter) > 0 || len(toIFilter) > 0
}

// corresponds to o_i_probs = 1 - np.prod(1 - (M * disease.trans_prob)[i_filter], axis=0)
func calculateToIProbs(M *mat.Dense, disease Disease, iFilter []int) []float64 {
	N, _ := M.Dims()
	probOfNoTransMatrix := mat.NewDense(N, N, nil)
	// (M * disease.trans_prob)
	probOfNoTransMatrix.Apply(func(i, j int, v float64) float64 {
		return v * disease.TransProb
	}, M)
	// (M * disease.trans_prob)[i_filter]
	probOfNoTransMatrix = newMatrixFromRows(probOfNoTransMatrix, iFilter)
	// 1 - (M * disease.trans_prob)[i_filter]
	probOfNoTransMatrix.Apply(func(i, j int, v float64) float64 { return 1 - v },
		probOfNoTransMatrix)
	// np.prod(1 - (M * disease.trans_prob)[i_filter], axis=0)
	nodeToTransProb := colProd(probOfNoTransMatrix)
	// 1 - np.prod(1 - (M * disease.trans_prob)[i_filter], axis=0)
	for i, v := range nodeToTransProb {
		nodeToTransProb[i] = 1 - v
	}
	return nodeToTransProb
}

func newMatrixFromRows(m *mat.Dense, rows []int) *mat.Dense {
	N, _ := m.Dims()
	nRows := len(rows)
	backingSlice := make([]float64, N*nRows)
	for i, rowNum := range rows {
		row := m.RawRowView(rowNum)
		for j, v := range row {
			backingSlice[i*N+j] = v
		}
	}
	return mat.NewDense(nRows, N, backingSlice)
}

// return the product of each of the columns in the matrix
func colProd(m *mat.Dense) []float64 {
	r, _ := m.Dims()
	matrixSlice := make([][]float64, r)
	for i := range matrixSlice {
		matrixSlice[i] = m.RawRowView(i)
	}
	var prod []float64
	copy(prod, matrixSlice[0])
	for _, row := range matrixSlice[1:] {
		for i, v := range row {
			prod[i] *= v
		}
	}
	return prod
}

func makeToIFilter(sir SIR, toIProbs []float64, rng distuv.Rander) []int {
	susceptibleAgents := sir.SusceptibleAgents()
	toIFilter := make([]int, 0)
	for agent := range susceptibleAgents {
		if rng.Rand() < toIProbs[agent] {
			toIFilter = append(toIFilter, agent)
		}
	}
	return toIFilter
}