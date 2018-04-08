package graph

import "utils"
import "cache"

type Graph interface {
	PlainReport()
	JsonReport()
	CacheAlgorithm() string
	LibrarySize() int
	SpectrumManager() SpectrumManager
	MatchStoreData(utils.BestSeparatorRanks, int) bool
	GenerateBestSeparatorRanksData([]int) utils.BestSeparatorRanks
	Clients() []Client
	ResetCounters()
	Links() []UnidirectionalLink
	OriginServers() []Node
	SpectrumCapacity() int
}

type Client interface {
	RandomRequest() interface{}
	RequestByID(int) interface{}
}

type SpectrumManager interface {
	BestSeparatorRanks([]int) []int
	BestReferenceRanks(utils.MirageStore) []int
	SetContentSpectrums([]int)
}

type ServerModel interface {
	Storage() cache.Storage
	AcceptRequest(cache.ContentRequest) interface{}
	ID() string
}

type UnidirectionalLink interface {
	Cost() float64
	Src() Node
	Dst() Node
	Traffic() float64
	SetTraffic(float64)
}

type Node interface {
	SetDijkstraCost(float64)
	SetDijkstraDone(bool)
	DijkstraCost() float64
	DijkstraDone() bool
	Entity() interface{}
	ID() string
	OutputLinks() []UnidirectionalLink
}

type Router interface {
	SelectForwardEntryBySpectrum(string, cache.ContentRequest) ForwardEntry
	SelectForwardEntry(string, cache.ContentRequest) ForwardEntry
	Init()
	ForwardRequest(string, cache.ContentRequest) interface{}
	RoutingTable() map[string]map[string][]ForwardEntry
}
