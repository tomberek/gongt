//
// Copyright (C) 2017 Yahoo Japan Corporation:
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package main

import (
	"bufio"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/kpango/glg"
	"github.com/yahoojapan/gongt"
)

const SIZEOF_FLOAT64 = 8

func getVectorsChanBinary(path string, vecChan chan<- []float64) ([][]float64, error) {
	fd, err := os.Open(path)
	defer close(vecChan)
	if err != nil {
		return nil, err
	}
	s, err := fd.Stat()
	if err != nil {
		return nil, err
	}
	_ = s
	var floatSlice [1024]float64
	for {
		err = binary.Read(fd, binary.LittleEndian, &floatSlice)
		if err == io.EOF {
			err = nil
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		vecChan <- floatSlice[:]
	}
	if err != nil {
		return nil, err
	}
	// Don't accumulate, TODO: remove return
	var res [][]float64
	return res, nil
}
func getVectorsChan(path string, vecChan chan<- []float64) ([][]float64, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	r := bufio.NewReader(f)
	var result [][]float64
	for {
		line, err := r.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		floats, err := parseFloats(line)
		if err != nil {
			log.Fatal(err)
		}
		vecChan <- floats
		result = append(result, floats)
	}
	close(vecChan)
	return result, nil
}
func getVectors(path string) ([][]float64, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	r := bufio.NewReader(f)
	var result [][]float64
	for {
		line, err := r.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		floats, err := parseFloats(line)
		if err != nil {
			log.Fatal(err)
		}
		result = append(result, floats)
	}
	return result, nil
}
func parseFloats(s string) ([]float64, error) {
	var (
		fields = strings.Fields(s)
		floats = make([]float64, len(fields))
		err    error
	)
	for i, f := range fields {
		floats[i], err = strconv.ParseFloat(f, 64)
		if err != nil {
			return nil, err
		}
	}
	return floats, nil
}

func create(name, path string, binary bool) {
	if _, err := os.Stat(name); err == nil {
		glg.Infof("[%s] %s exists", name, path)
		return
	}
	vecChan := make(chan []float64)
	done := make(chan struct{})
	go func(vecChan <-chan []float64) {
		v := <-vecChan
		n := gongt.New("a").SetObjectType(gongt.Float).
			SetDimension(len(v)).
			SetBulkInsertChunkSize(200).
			SetCreationEdgeSize(20).
			SetSearchEdgeSize(40).
			SetDistanceType(gongt.Cosine).
			Open()

		defer n.Close()
		var wg sync.WaitGroup

		n.Insert(v)
		var i = 1
		for v := range vecChan {
			n.Insert(v)
			i = i + 1
			if i%100000 == 0 {
				glg.Infof("Processing: %s objects", i)
				// wg.Add(1)
				// go func(wg *sync.WaitGroup) {
				// 	defer wg.Done()
				// if err := n.CreateAndSaveIndex(runtime.NumCPU()); err != nil {
				// 	glg.Warn(err)
				// }
				// glg.Infof("Saved: %s objects", i)
				// }(&wg)
			}
		}
		wg.Wait()
		if err := n.CreateAndSaveIndex(runtime.NumCPU()); err != nil {
			glg.Warn(err)
		}
		glg.Info("Finished with indexing")
		close(done)
	}(vecChan)

	var err error
	if binary {
		_, err = getVectorsChanBinary(path, vecChan)
	} else {
		_, err = getVectorsChan(path, vecChan)
	}
	if err != nil {
		glg.Warn(err)
		return
	}
	<-done
	defer glg.Infof("done")

	// _, errs := n.BulkInsert(vectors)
	// for _, err := range errs {
	// 	if err != nil {
	// 		glg.Warn(err)
	// 		return
	// 	}
	// }
	// if err := n.CreateAndSaveIndex(runtime.NumCPU()); err != nil {
	// if err := n.CreateAndSaveIndex(16); err != nil {
	// 	glg.Warn(err)
	// }
}

func search(name, path string) {
	n := gongt.New(name).Open()
	defer n.Close()

	vectors, err := getVectors(path)
	if err != nil {
		glg.Warn(err)
		return
	}
	glg.Infof("[%s] %d items", name, len(vectors))
	defer glg.Infof("[%s] done", name)

	for _, v := range vectors {
		n.Search(v, 10, gongt.DefaultEpsilon)
		// result, err := n.Search(v, 10, gongt.DefaultEpsilon) // do something using result and err
	}
}

func extract(name string) {
	n := gongt.New(name).Open()
	defer n.Close()

	res, err := n.ExtractGraph()
	if err != nil {
		glg.Warn(err)
		return
	}
	for _, item := range res {
		var line []string
		for _, s := range item {
			line = append(line, fmt.Sprintf("%d %f", s.ID, s.Distance))
		}
		output := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(line)), "\t"), "[]")
		fmt.Println(output)
	}
}

func main() {
	c := flag.Bool("create", false, "run create")
	b := flag.Bool("binary", false, "create with binary format")
	s := flag.Bool("search", false, "run search")
	e := flag.Bool("extract", false, "run extract")

	n := flag.String("name", "", "dataset name")
	p := flag.String("path", "", "dataset path")

	flag.Parse()
	if *c {
		create(*n, *p, *b)
	}
	if *s {
		search(*n, *p)
	}
	if *e {
		extract(*n)
	}
}
