package main

import (
	"distribution/userdist"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var filename string = "test_data.csv"

func getDataDist() map[int]float64 {
	contents, err := ioutil.ReadFile(filename)
	data_map := make(map[int]float64)

	if err != nil {
		panic(fmt.Errorf("error in reading file: %s", err))
	}
	lines := strings.Split(string(contents), "\n")
	for index := range lines {
		line := lines[index]
		popularity, err := strconv.ParseFloat(line, 64)
		if err != nil {
			continue
		}
		//dist.rawPopularities = append(dist.rawPopularities, algorithm.Entry{line, -popularity})
		data_map[index+1] = popularity
	}
	return data_map
}

func main() {
	dist := userdist.New("test_data.csv")
	dic := make(map[int]int)
	Requests := 1500000
	for i := 0; i < Requests; i++ {
		dic[dist.Intn()]++
	}
	fmt.Println(dic)

	data_map := getDataDist()
	fmt.Println(data_map)

	var sum float64 = 0.0
	for i := 0; i < len(data_map); i++ {
		sum += data_map[i+1]
	}
	fmt.Println(sum)

	NoC := len(dic)
	fmt.Printf("%s \t\t %s \t\t %s \t %s\n",
		"Rank", "Count", "Real Probability", "Generated Probability")

	for i := 0; i < NoC; i++ {
		fmt.Printf("%d \t\t %d \t\t %f \t\t %f\n",
			i+1, dic[i+1], data_map[i+1]/sum, float64(dic[i+1])/float64(Requests))
	}
}
