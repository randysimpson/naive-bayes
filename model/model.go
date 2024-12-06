/*
MIT License

# Copyright (©) 2024 - Randall Simpson

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/
package model

import (
	"math"

	"k8s.io/klog"
)

var count map[string]int
var bicount map[string]map[string]int
var tricount map[string]map[string]map[string]int
var quadcount map[string]map[string]map[string]map[string]int

// laplase smoothing
var laplace_alpha float64

type Key struct {
	First  string `json:"first"`
	Second string `json:"second"`
}

var bigramModel map[Key]float64

// trigram
type TriKey struct {
	First  string `json:"first"`
	Second string `json:"second"`
	Third  string `json:"thrid"`
}

var trigramModel map[TriKey]float64

// quadgram
type QuadKey struct {
	First  string `json:"first"`
	Second string `json:"second"`
	Third  string `json:"thrid"`
	Fourth string `json:"fourth"`
}

var quadgramModel map[QuadKey]float64

func init() {
	ClearModel()
}

func ClearModel() {
	count = map[string]int{}
	bicount = map[string]map[string]int{}
	tricount = map[string]map[string]map[string]int{}
	quadcount = map[string]map[string]map[string]map[string]int{}
	laplace_alpha = 0.001

	bigramModel = map[Key]float64{}
	trigramModel = map[TriKey]float64{}
	quadgramModel = map[QuadKey]float64{}
}

type CountModel struct {
	Count     map[string]int                                  `json:"count"`
	BiCount   map[string]map[string]int                       `json:"biCount"`
	TriCount  map[string]map[string]map[string]int            `json:"triCount"`
	QuadCount map[string]map[string]map[string]map[string]int `json:"quadCount"`
}

type KeyModels struct {
	BiModel   map[Key]float64     `json:"biModel"`
	TriModel  map[TriKey]float64  `json:"triModel"`
	QuadModel map[QuadKey]float64 `json:"quadModel"`
}

type RtnModel struct {
	CountModels CountModel `json:"countModels"`
	KeyModels   KeyModels  `json:"keyModels"`
}

func GetModel() RtnModel {
	return RtnModel{
		CountModels: CountModel{
			Count:     count,
			BiCount:   bicount,
			TriCount:  tricount,
			QuadCount: quadcount,
		},
		KeyModels: KeyModels{
			BiModel:   bigramModel,
			TriModel:  trigramModel,
			QuadModel: quadgramModel,
		},
	}
}

func addBigram(first string, second string) {
	mm, ok := bicount[first]
	if !ok {
		mm = map[string]int{}
		bicount[first] = mm
	}
	mm[second]++
}

func addTrigram(first string, second string, third string) {
	_, ok := tricount[first]
	if !ok {
		tricount[first] = map[string]map[string]int{}
	}
	m2, ok := tricount[first][second]
	if !ok {
		m2 = map[string]int{}
		tricount[first][second] = m2
	}
	m2[third]++
}

func addQuadgram(first string, second string, third string, fourth string) {
	_, ok := quadcount[first]
	if !ok {
		quadcount[first] = map[string]map[string]map[string]int{}
	}
	_, ok = quadcount[first][second]
	if !ok {
		quadcount[first][second] = map[string]map[string]int{}
	}
	m3, ok := quadcount[first][second][third]
	if !ok {
		m3 = map[string]int{}
		quadcount[first][second][third] = m3
	}
	m3[fourth]++
}

func AddData(data []string) (int, error) {
	for i, word := range data {
		count[word]++
		if i < len(data)-1 {
			addBigram(word, data[i+1])
		}
		if i < len(data)-2 {
			addTrigram(word, data[i+1], data[i+2])
		}
		if i < len(data)-3 {
			addQuadgram(word, data[i+1], data[i+2], data[i+3])
		}
	}

	klog.Infof("count: %+v", count)
	klog.Infof("bicount: %+v", bicount)
	klog.Infof("tricount: %+v", tricount)
	klog.Infof("quadcount: %+v", quadcount)

	buildModel()

	return len(data), nil
}

func GetCount() int {
	return len(count)
}

func GetNext(word string) (map[string]float32, error) {
	//find all bi-counts with words
	probabilities := map[string]float32{}

	total := float32(count[word])

	for option, value := range bicount[word] {
		probabilities[option] = float32(value) / total
	}

	return probabilities, nil
}

func GetTriNext(first string, second string) (map[string]float32, error) {
	probabilities := map[string]float32{}

	total := 0.0

	for option, value := range tricount[first][second] {
		total += float64(value)
		probabilities[option] = float32(value)
	}

	for option, value := range probabilities {
		probabilities[option] = float32(value) / float32(total)
	}

	return probabilities, nil
}

func GetQuadNext(first string, second string, third string) (map[string]float32, error) {
	probabilities := map[string]float32{}

	total := 0.0

	for option, value := range quadcount[first][second][third] {
		total += float64(value)
		probabilities[option] = float32(value)
	}

	for option, value := range probabilities {
		probabilities[option] = float32(value) / float32(total)
	}

	return probabilities, nil
}

func GetEntropy(data []string) (float64, error) {
	n := len(data)

	var total float64
	total = 0.0

	for i, word := range data {
		if i < len(data)-1 {
			total += getProbability(word, data[i+1])
		}
	}

	exponent := total * (float64(-1) / float64(n))
	perplexity := math.Pow(250, exponent)

	return perplexity, nil
}

func getProbability(first string, second string) float64 {
	//check if word exists
	if count[first] == 0 {
		first = "<UKN>"
	}
	if count[second] == 0 {
		second = "<UKN>"
	}
	return bigramModel[Key{first, second}]
}

func buildModel() {
	//use log base of 250
	logBase := 1 / math.Log(250)

	v := float64(len(count)) + laplace_alpha

	for key, value := range count {
		denom := float64(value) + v
		for key2, val2 := range bicount[key] {
			bigramModel[Key{key, key2}] = (float64(val2) + laplace_alpha) / denom

			//trigram
			for key3, val3 := range tricount[key][key2] {
				trigramModel[TriKey{key, key2, key3}] = (float64(val3) + laplace_alpha) / denom

				//quadgram
				for key4, val4 := range quadcount[key][key2][key3] {
					quadgramModel[QuadKey{key, key2, key3, key4}] = (float64(val4) + laplace_alpha) / denom
				}
				//unknown token
				quadgramModel[QuadKey{key, key2, key3, "<UKN>"}] = math.Log(laplace_alpha/denom) * logBase
			}
			//unknown token
			trigramModel[TriKey{key, key2, "<UKN>"}] = math.Log(laplace_alpha/denom) * logBase
		}
		//add unknown token
		bigramModel[Key{key, "<UKN>"}] = math.Log(laplace_alpha/denom) * logBase
		trigramModel[TriKey{key, "<UKN>", "<UKN>"}] = math.Log(laplace_alpha/denom) * logBase
	}
	//handle unknown as first word.
	bigramModel[Key{"<UKN>", "<UKN>"}] = math.Log(laplace_alpha/(v+laplace_alpha)) * logBase
	trigramModel[TriKey{"<UKN>", "<UKN>", "<UKN>"}] = math.Log(laplace_alpha/(v+laplace_alpha)) * logBase
}
