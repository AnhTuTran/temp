package graph

import "fmt"
import "utils"
import "algorithm"

type TrafficExpectation_t struct {
	network Graph
	Traffic map[UnidirectionalLink]float64
}

func (expectation *TrafficExpectation_t) SetNetwork(network Graph) {
	expectation.network = network
}

func (expectation *TrafficExpectation_t) TotalTraffic() float64 {
	totalTraffic := 0.0
	for _, link := range expectation.network.Links() {
		totalTraffic += expectation.Traffic[link]
	}
	return totalTraffic
}

func (expectation *TrafficExpectation_t) originTraffic() float64 {
	originTraffic := 0.0
	for _, origin := range expectation.network.OriginServers() {
		for _, link := range origin.OutputLinks() {
			originTraffic += expectation.Traffic[link]
		}
	}
	return originTraffic
}

func (expectation *TrafficExpectation_t) internalTraffic() float64 {
	return expectation.TotalTraffic() - expectation.originTraffic()
}

func (expectation *TrafficExpectation_t) inspect() {
	utils.DebugPrint(fmt.Sprintln("--"))
	for _, link := range expectation.network.Links() {
		utils.DebugPrint(fmt.Sprintf("%s -> %s :\t%f\n", link.Src().ID, link.Dst().ID, expectation.Traffic[link]))
	}
	utils.DebugPrint(fmt.Sprintln("--"))
	linkTraffic := make([]float64, 0)
	for _, traffic := range expectation.Traffic {
		linkTraffic = append(linkTraffic, traffic)
	}
	utils.DebugPrint(fmt.Sprintf("Link traffic SD: %f\n", algorithm.StandardDeviation(linkTraffic)))
	utils.DebugPrint(fmt.Sprintln("--"))
	utils.DebugPrint(fmt.Sprintln("total_traffic: ", expectation.TotalTraffic()))
	utils.DebugPrint(fmt.Sprintln(" - internal_traffic:", expectation.internalTraffic()))
	utils.DebugPrint(fmt.Sprintln(" - origin_traffic:  ", expectation.originTraffic()))
}
