package main

import "utils"

type EstimationResult struct {
	separatorRanks []int
	trafficSize    float64
}

type SpectrumManager_t struct {
	network              *Graph_t
	baseSpectrums        []interface{}
	spectrumTags         map[int][]uint64
	bitSize              int
	contentSpectrums     []uint64
	serverSpectrums      map[*Node_t]uint64
	spectrumRoutingTable map[*Node_t]map[uint64][]ForwardEntry
}

func newSpectrumManager(bitSize int, network *Graph_t) *SpectrumManager_t {
	manager := new(SpectrumManager_t)
	return manager
}

func (manager *SpectrumManager_t) initSpectrumRoutingTable() {
}

func (manager *SpectrumManager_t) inspectSpectrumRoutingTable() {
}

func (manager *SpectrumManager_t) adjacentSpectrums(node *Node_t) []uint64 {
	adjacentSpectrums := make([]uint64, 0)
	return adjacentSpectrums
}

func (manager *SpectrumManager_t) availableSpectrums(node *Node_t) []uint64 {
	availableSpectrums := make([]uint64, 0)
	return availableSpectrums
}

func (manager *SpectrumManager_t) countSpectrums(spectrum uint64) int {
	return 0
}

func (manager *SpectrumManager_t) selectDistantSpectrum(srcNode *Node_t, availableSpectrums []uint64) uint64 {
	return 0
}

func (manager *SpectrumManager_t) setServerSpectrums() {
}

func (manager *SpectrumManager_t) inspectServerSpectrums() {
}

func (manager *SpectrumManager_t) SetContentSpectrums(separatorRanks []int) {
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
