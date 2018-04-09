package main

import (
	"fmt"
	"math"
	"reflect"
)

// TU ALL

type SeparatorRanks_t struct {
	network *Graph_t
}

func isTwoArraysDiff(a, b []int) bool {
	len := len(a)

	for i := 0; i < len; i++ {
		if a[i] != b[i] {
			return true
		}
	}
	return false
}

func getValues(a []int) []int {
	len := len(a)
	out := make([]int, len)
	for i := 0; i < len; i++ {
		out[i] = a[i]
	}
	return out
}

func max(a, b int) int {
	if a >= b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

func calculateTail(S []int, N, C int) int {
	result := C * N
	end_i := len(S) - 1
	for i := 0; i < end_i; i++ {
		if i == 0 {
			result -= S[i] * N
		} else {
			result -= (S[i] - S[i-1]) * (N - i)
		}
	}
	return result
}

func (separatorranks *SeparatorRanks_t) estimate_traffic(S_tmp []int,
	numUsrs, numServers, numContents int) float64 {

	traffic := 0.0
	routingTable := separatorranks.network.router.(*Router_t).routingTable
	graph := separatorranks.network

	fmt.Println("Empty")
	fmt.Println(routingTable["node1"]["node1"])
	// src := graph.clients[1].Upstream().ID()
	// des := graph.nodes[0].id
	// fmt.Println(src)
	// fmt.Println(graph.nodes)

	dist := graph.clients[0].Dist()

	for i := 0; i < numUsrs; i++ {
		for j := 0; j < numServers; j++ {
			fmt.Printf("i %d j %d\n", i, j)
			src := graph.clients[i].Upstream().ID()
			des := graph.nodes[j].id
			tmp := routingTable[src][des]
			cost := 0.0
			if len(tmp) == 0 {
				cost = 1.0
			} else {
				cost = tmp[0].Cost()
			}
			for k := 0; k < numContents; k++ {
				//nearest_server := src
				bin_var := 0.0
				if graph.clients[i].Upstream().Storage().Exist(k + 1) {
					bin_var = 1.0
				}

				traffic += cost * dist.PDF(k+1) * bin_var
				// fmt.Printf("k %d\n", k)
			}
		}
	}
	fmt.Println(traffic)
	return traffic
}

func newSeparatorRanks(network *Graph_t) *SeparatorRanks_t {
	separator := new(SeparatorRanks_t)
	separator.network = network
	return separator
}

func (separatorranks *SeparatorRanks_t) GetSeparatorRanks() {
	fmt.Println("XXXX")
	//fmt.Println(separatorranks.network.router.(*Router_t).routingTable)
	fmt.Println("XXXX")

	//routingTable := separatorranks.network.router.(*Router_t).routingTable
	separatorranks.network.router.(*Router_t).inspectAllRoutingTables()

	// for i := range routingTable {
	// 	fmt.Println(routingTable[i])
	// }

	// for nextNode := routingTable["node3"]["origin"][0].Node().(*Node_t).id; nextNode != "origin"; {
	// 	fmt.Println(routingTable[nextNode]["origin"][0].Node())
	// 	nextNode = routingTable["node3"]["origin"][0].Node().(*Node_t).id
	// }
	// nextNode := routingTable["node3"]["origin"][0].Node().(*Node_t).id
	// fmt.Println(routingTable[nextNode]["origin"][0].Node())

	// nextNode = routingTable["node3"]["origin"][0].Node().(*Node_t).id
	// fmt.Println(routingTable[nextNode]["origin"][0].Node())

	N := 4   // # colors
	C := 100 // cache server capacity
	numUsrs := len(separatorranks.network.clients)
	numServers := len(separatorranks.network.nodes)
	numContents := separatorranks.network.LibrarySize()

	fmt.Println("XXXX")
	fmt.Println(numUsrs)
	fmt.Println(numServers)
	fmt.Println(numContents)

	S := make([]int, 4)
	S_prev := make([]int, 4)
	var S_tmp []int
	T_min := math.MaxFloat64

	fmt.Println(reflect.TypeOf(T_min))

	S[N-1] = N * C
	fmt.Println(S)

	for isTwoArraysDiff(S, S_prev) {
		S_prev = getValues(S)
		for i := 0; i <= N-2; i++ {
			start_v := max(0, S[max(1, i)-1])
			end_v := min(S[i+1], N*C)
			for v := start_v; v <= end_v; v++ {
				S_tmp = getValues(S)
				S_tmp[i] = v
				S_tmp[N-1] = calculateTail(S_tmp, N, C)
			}
			T_est := separatorranks.estimate_traffic(S_tmp, numUsrs, numServers, numContents)
			if T_est < T_min {
				T_min = T_est
				S = getValues(S_tmp)
			}
		}
	}

	fmt.Println(S)
}
