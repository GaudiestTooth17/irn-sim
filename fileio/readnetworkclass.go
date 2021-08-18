package fileio

import (
	"archive/tar"
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/GaudiestTooth17/irn-sim/network"
)

func ReadClass(pathToClass string) []*network.AdjacencyList {
	reader, err := os.Open(pathToClass)
	if err != nil {
		panic(err)
	}
	defer reader.Close()
	tarReader := tar.NewReader(reader)
}

func untar(tarball, target string) {
	reader, err := os.Open(tarball)
	if err != nil {
		panic(err)
	}
	defer reader.Close()
	tarReader := tar.NewReader(reader)

	for {
		header, err := tarReader.Next()
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			panic(err)
		}

		path := filepath.Join(target, header.Name)
		info := header.FileInfo()

		file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
		if err != nil {
			panic(err)
		}
		defer file.Close()
		_, err = io.Copy(file, tarReader)
	}
}
