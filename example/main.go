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
	"flag"
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

func create(name, path string) {
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
			SetCreationEdgeSize(10).
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
				wg.Add(1)
				go func(wg *sync.WaitGroup) {
					defer wg.Done()
					if err := n.CreateAndSaveIndex(30); err != nil {
						glg.Warn(err)
					}
					glg.Infof("Saved: %s objects", i)
				}(&wg)
			}
		}
		wg.Wait()
		if err := n.CreateAndSaveIndex(30); err != nil {
			glg.Warn(err)
		}
		glg.Info("Finished with indexing")
		close(done)
	}(vecChan)

	_, err := getVectorsChan(path, vecChan)
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
	_ = runtime.NumCPU()
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

func main() {
	c := flag.Bool("create", false, "run create")
	s := flag.Bool("search", false, "run search")

	n := flag.String("name", "", "dataset name")
	p := flag.String("path", "", "dataset path")

	flag.Parse()
	if *c {
		create(*n, *p)
	}
	if *s {
		search(*n, *p)
	}
}
