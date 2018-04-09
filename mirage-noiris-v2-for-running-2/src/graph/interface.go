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
	//TU
	SeparatorRanks() SeparatorRanks
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

//TU
type SeparatorRanks interface {
	GetSeparatorRanks()
}

type ServerModel interface {
	Storage() cache.Storage
	AcceptRequest(cache.ContentRequest) interface{}
	ID() string
}

type UnidirectionalLink interface {
	Src() Node
	Dst() Node
	Traffic() float64
	SetTraffic(float64)
}

type Node interface {
	Entity() interface{}
	ID() string
	OutputLinks() []UnidirectionalLink
}
