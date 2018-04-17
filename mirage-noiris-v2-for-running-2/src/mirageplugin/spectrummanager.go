package main

import (
	"container/list"
	"fmt"
	"sort"
	"utils"
)

type EstimationResult struct {
	separatorRanks []int
	trafficSize    float64
}

type SpectrumManager_t struct {
	network              *Graph_t
	baseSpectrums        []interface{}                         // colors not sure ???
	spectrumTags         map[int][]uint64                      // contents' color
	bitSize              int                                   // # colors
	contentSpectrums     []uint64                              // seperatorRanks
	serverSpectrums      map[*Node_t]uint64                    // servers' color
	spectrumRoutingTable map[*Node_t]map[uint64][]ForwardEntry // color-based routing table
}

func newSpectrumManager(bitSize int, network *Graph_t) *SpectrumManager_t {
	manager := new(SpectrumManager_t)
	// TU ADD
	manager.network = network
	manager.bitSize = bitSize
	manager.serverSpectrums = make(map[*Node_t]uint64)
	manager.baseSpectrums = make([]interface{}, 0)
	for i := 0; i < bitSize; i++ {
		manager.baseSpectrums = append(manager.baseSpectrums, uint64(i))
	}
	//
	return manager
}

func (manager *SpectrumManager_t) initSpectrumRoutingTable() {
}

func (manager *SpectrumManager_t) inspectSpectrumRoutingTable() {
}

func (manager *SpectrumManager_t) initSpectrums() {
}

func (manager *SpectrumManager_t) initBaseSpectrums() {
}

func (manager *SpectrumManager_t) estimateTotalTraffic(separatorRanks []int, join chan EstimationResult) {
}

func (manager *SpectrumManager_t) adjustSeparatorRanks(separatorRanks []int, librarySize int) bool {
	return false
}

func (manager *SpectrumManager_t) separatorRanksID(separatorRanks []int) string {
	return ""
}

func (manager *SpectrumManager_t) BestSeparatorRanks(separatorRanks []int) []int {
	return make([]int, 0)
}

func (manager *SpectrumManager_t) BestReferenceRanks(mirageStore utils.MirageStore) []int {
	return make([]int, 0)
}

// TU modified code
func isInArray(array []uint64, elem uint64) bool {
	exist := false
	for _, v := range array {
		if elem == v {
			exist = true
			break
		}
	}
	return exist
}

func (manager *SpectrumManager_t) adjacentSpectrums(node *Node_t) []uint64 {
	adjacentSpectrums := make([]uint64, 0)
	for _, link := range node.inputLinks {
		color, ok := manager.serverSpectrums[link.src]
		if ok {
			if !isInArray(adjacentSpectrums, color) {
				adjacentSpectrums = append(adjacentSpectrums, color)
			}
		}
	}
	for _, link := range node.outputLinks {
		color, ok := manager.serverSpectrums[link.dst]
		if ok {
			if !isInArray(adjacentSpectrums, color) {
				adjacentSpectrums = append(adjacentSpectrums, color)
			}
		}
	}
	return adjacentSpectrums
}

func (manager *SpectrumManager_t) availableSpectrums(node *Node_t) []uint64 {
	availableSpectrums := make([]uint64, 0)
	adjacentSpectrums := manager.adjacentSpectrums(node)
	for _, v := range manager.baseSpectrums {
		if !isInArray(adjacentSpectrums, v.(uint64)) {
			availableSpectrums = append(availableSpectrums, v.(uint64))
		}
	}
	return availableSpectrums
}

func (manager *SpectrumManager_t) countSpectrums(spectrum []uint64) int {
	return len(spectrum)
}

func (manager *SpectrumManager_t) selectDistantSpectrum(srcNode *Node_t, availableSpectrums []uint64) uint64 {
	// array stores minimal distance to node with available colors
	distances := make([]int, len(availableSpectrums))
	for i := range distances {
		distances[i] = -1 // distance value not found
	}
	// list to do BFS
	list := list.New()
	hopcount := 1

	for i := 0; i < len(srcNode.outputLinks); i++ {
		list.PushBack(srcNode.outputLinks[i].dst)
		list.PushBack(hopcount)
	}

	// fmt.Println(ele.Value)
	for !isAllElePositive(distances) {
		isInAvaiSpectrums := true

		node := list.Remove(list.Front()).(*Node_t)
		//fmt.Println(node.id)
		distance := list.Remove(list.Front()).(int)
		nodeColor, exist := manager.serverSpectrums[node]
		if exist {
			for i := 0; i < len(availableSpectrums); i++ {
				if nodeColor == availableSpectrums[i] {
					distances[i] = distance
					break
				}
			}
			isInAvaiSpectrums = false
		}

		if !exist || !isInAvaiSpectrums {
			hopcount++
			for i := 0; i < len(node.outputLinks); i++ {
				list.PushBack(node.outputLinks[i].dst)
				list.PushBack(hopcount)
			}
		}

	}
	minDistance := distances[0]
	selectedColor := availableSpectrums[0]
	for i := 0; i < len(availableSpectrums); i++ {
		if distances[i] < minDistance {
			minDistance = distances[i]
			selectedColor = availableSpectrums[i]
		}
	}
	return selectedColor
}

func (manager *SpectrumManager_t) setServerSpectrums() {
	verticesDegrees := make([]vertexDegree, 0)
	for _, node := range manager.network.nodes {
		verticesDegrees = append(verticesDegrees,
			vertexDegree{node, uint64(len(node.inputLinks) + len(node.outputLinks))})
	}

	sort.Slice(verticesDegrees, func(i, j int) bool {
		return verticesDegrees[i].degree > verticesDegrees[j].degree
	})

	for i := range verticesDegrees {
		availableSpectrums := manager.availableSpectrums(verticesDegrees[i].node)
		// missing sort descendingly based on minimal distance
		selectedColor := uint64(manager.bitSize)
		if i >= manager.bitSize {
			selectedColor = manager.selectDistantSpectrum(verticesDegrees[i].node, availableSpectrums)
		} else {
			selectedColor = manager.baseSpectrums[i].(uint64)
		}
		manager.serverSpectrums[verticesDegrees[i].node] = selectedColor
	}
	fmt.Println("\nServers' color")
	for _, node := range manager.network.nodes {
		fmt.Printf("%s %d\n", node.id, manager.serverSpectrums[node])
	}
}

func (manager *SpectrumManager_t) inspectServerSpectrums() {
}

func (manager *SpectrumManager_t) SetContentSpectrums(separatorRanks []int) {
}

// TU add code
type vertexDegree struct {
	node   *Node_t
	degree uint64
}

// type listElement {
// 	node *Node_t
// 	distance int
// }

func isAllElePositive(array []int) bool {
	for _, v := range array {
		if v <= 0 {
			return false
		}
	}
	return true
}
