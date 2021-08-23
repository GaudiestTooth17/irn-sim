package fileio

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"

	"github.com/GaudiestTooth17/irn-sim/network"
)

var idMatcher *regexp.Regexp = regexp.MustCompile(`\d+`)

// pathToClass must point to a .tar.gz file
func ReadClass(pathToClass string) []*network.AdjacencyList {
	// extract the compressed networks into /tmp (which is hopefully stored in RAM)
	className := filepath.Base(pathToClass)
	//strip extension
	className = className[:len(className)-7]
	extractionDest := filepath.Join("/tmp", className)
	if _, err := os.Stat(extractionDest); os.IsNotExist(err) {
		os.Mkdir(extractionDest, fs.ModePerm)
		ungzipAndUntar(extractionDest, pathToClass)
	}

	// collect all instances of networks=
	classInstances := make([]string, 0)
	filepath.WalkDir(extractionDest, func(path string, d fs.DirEntry, err error) error {
		filename := filepath.Base(path)
		isClassInstance, _ := regexp.Match(`instance-\d+.txt`, []byte(filename))
		if fileInfo, _ := os.Stat(path); !fileInfo.IsDir() && isClassInstance {
			classInstances = append(classInstances, path)
		}
		return nil
	})
	// sort according to the ID's at the end of the file name
	sort.Slice(classInstances, func(i, j int) bool {
		nameI := filepath.Base(classInstances[i])
		nameJ := filepath.Base(classInstances[j])
		idI := atoiOrPanic(idMatcher.FindString(nameI))
		idJ := atoiOrPanic(idMatcher.FindString(nameJ))
		return idI < idJ
	})

	// load the networks into memory
	nets := make([]*network.AdjacencyList, len(classInstances))
	for i, pathToInstance := range classInstances {
		nets[i] = ReadFile(pathToInstance)
	}
	return nets
}

func ungzipAndUntar(target, gzippedTarball string) {
	file, err := os.Open(gzippedTarball)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		panic(err)
	}
	defer gzipReader.Close()

	tarReader := tar.NewReader(gzipReader)

	for {
		header, err := tarReader.Next()
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			panic(err)
		}

		path := filepath.Join(target, filepath.Base(header.Name))
		info := header.FileInfo()
		if info.IsDir() {
			continue
		}

		file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
		if err != nil {
			panic(err)
		}
		defer file.Close()
		_, err = io.Copy(file, tarReader)
		if err != nil {
			panic(err)
		}
	}
}

func atoiOrPanic(str string) int {
	integer, err := strconv.Atoi(str)
	if err != nil {
		panic(err)
	}
	return integer
}
