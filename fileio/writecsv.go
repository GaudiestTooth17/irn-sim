package fileio

import (
	"encoding/csv"
	"os"
)

func WriteToCSV(csvName string, lines [][]string) {
	outFile, err := os.Create(csvName)
	if err != nil {
		panic(err)
	}
	defer outFile.Close()
	outWriter := csv.NewWriter(outFile)

	for _, line := range lines {
		outWriter.Write(line)
	}
	outWriter.Flush()
}
