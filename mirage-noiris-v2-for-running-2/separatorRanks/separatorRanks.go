package main

import (
	"fmt"
	"math"
	"reflect"
)

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

func estimate_traffic(S_tmp []int) int {
	traffic := 0
	return traffic
}

func main() {
	// if parser.Options.GraphFilename == "" {
	// 	utils.DebugPrint(fmt.Sprintln("Please specify graph filename."))
	// 	parser.PrintDefaults()
	// 	os.Exit(1)
	// }

	// network, _ := graphLoader.LoadGraph(parser.Options.GraphFilename)
	// warmupRequestCount := parser.Options.WarmupRequestCount
	// evaluationRequestCount := parser.Options.EvaluationRequestCount

	// generate requests for warming up
	// for index := 0; index < warmupRequestCount; index++ {
	// 	utils.DebugPrint(fmt.Sprintf("\rwarming up: (%d/%d)", index+1, warmupRequestCount))
	// 	for _, client := range network.Clients() {
	// 		client.RandomRequest()
	// 	}
	// }

	N := 4   // # colors
	C := 100 // cache server capacity

	S := make([]int, 4)
	S_prev := make([]int, 4)
	var S_tmp []int
	T_min := math.MaxInt64

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
			T_est := estimate_traffic(S_tmp)
			if T_est < T_min {
				T_min = T_est
				S = getValues(S_tmp)
			}
		}
	}

	fmt.Println(S)
}
