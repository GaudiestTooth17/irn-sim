package main

import (
	"fmt"
	"math/rand"
	"path/filepath"
	"time"

	fio "github.com/GaudiestTooth17/irn-sim/fileio"
	"github.com/GaudiestTooth17/irn-sim/network"
	"github.com/GaudiestTooth17/irn-sim/sim"
)

func main() {
	// if len(os.Args) < 2 {
	// 	fmt.Printf("Usage: %s <network>\n", os.Args[0])
	// 	return
	// }
	simsPerClassInstance := 1
	seed := int64(69)
	classPaths := getClassPaths()
	csvLines := make([][]string, len(classPaths)*2)
	for i, classPath := range classPaths {
		startTime := time.Now()
		networkName := filepath.Base(classPath)
		// strip off .tar.gz
		networkName = networkName[:len(networkName)-7]
		fmt.Printf("Running %d simulations on %s. ", simsPerClassInstance, networkName)

		nets := fio.ReadClass(classPath)
		// set up the parameters
		disease := sim.Disease{DaysInfectious: 4, TransProb: .2}
		makeBehavior := func(net *network.AdjacencyList, rng *rand.Rand) sim.Behavior {
			return sim.NewSimplePressureBehavior(net, rng, 2, .25)
		}
		makeSIR0 := func(N int, numToInfect int, rng *rand.Rand) sim.SIR {
			return sim.MakeSir0(N, 1, rng)
		}

		// run a simulations
		survivalRates := sim.SimOnManyNetworksForSurvivalRate(nets, makeSIR0, disease,
			makeBehavior, 300, seed, simsPerClassInstance)
		csvLines[2*i] = []string{networkName}
		csvLines[2*i+1] = floatSliceToStrSlice(survivalRates)

		// report completion
		fmt.Printf("Done (%v).\n", time.Since(startTime))
	}
	// save to csv
	fio.WriteToCSV("results/survival rates (go).csv", csvLines)
}

func floatSliceToStrSlice(fSlice []float64) []string {
	strSlice := make([]string, len(fSlice))
	for i, val := range fSlice {
		strSlice[i] = fmt.Sprint(val)
	}
	return strSlice
}

func getClassPaths() []string {
	networkPaths := []string{
		"BarabasiAlbert(N=500,m=2).tar.gz",
		"BarabasiAlbert(N=500,m=2).tar.gz",
		"BarabasiAlbert(N=500,m=2).tar.gz",
		"ConnComm(N_comm=10,ib=(5, 10),num_comms=50,ob=(3, 6)).tar.gz",
		"ConnComm(N_comm=20,ib=(15, 20),num_comms=25,ob=(3, 6)).tar.gz",
		"ErdosRenyi(N=500,p=0.01).tar.gz",
		"ErdosRenyi(N=500,p=0.02).tar.gz",
		"ErdosRenyi(N=500,p=0.03).tar.gz",
		"WattsStrogatz(N=500,k=4,p=0.01).tar.gz",
		"WattsStrogatz(N=500,k=4,p=0.02).tar.gz",
		"WattsStrogatz(N=500,k=5,p=0.01).tar.gz",
	}
	for i, path := range networkPaths {
		networkPaths[i] = filepath.Join("networks", path)
	}
	return networkPaths
}
