package main

import "fmt"
import "encoding/json"
import "io/ioutil"
import "graph"
import "distribution"
import "distribution/gamma"
import "distribution/zipf"

type RequestDecodeModel struct {
	ID            string    `json:"Id"`
	Model         string    `json:"Model"`
	ParameterKeys []string  `json:"ParameterKeys"`
	Parameters    []float64 `json:"Parameters"`
}

type OriginDecodeModel struct {
	ID          string `json:"Id"`
	LibrarySize int    `json:"LibrarySize"`
}

type NodeDecodeModel struct {
	ID             string    `json:"Id"`
	CacheAlgorithm string    `json:"CacheAlgorithm"`
	ParameterKeys  []string  `json:"ParameterKeys"`
	Parameters     []float64 `json:"Parameters"`
}

type LinkDecodeModel struct {
	EdgeNodeIds   []string `json:"EdgeNodeIds"`
	Cost          float64  `json:"Cost"`
	Bidirectional bool     `json:"Bidirectional"`
}

type ClientDecodeModel struct {
	RequestModelID string `json:"RequestModelId"`
	UpstreamID     string `json:"UpstreamId"`
	TrafficWeight  float64
}

type DecodeModel struct {
	NetworkID      string               `json:"NetworkId"`
	RequestModels  []RequestDecodeModel `json:"RequestModels"`
	Origin         OriginDecodeModel    `json:"OriginServer"`
	Nodes          []NodeDecodeModel    `json:"CacheServers"`
	Links          []LinkDecodeModel    `json:"Links"`
	Clients        []ClientDecodeModel  `json:"Clients"`
	ModelRelations map[interface{}]interface{}
}

func loadFile(filename string) DecodeModel {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(fmt.Errorf("Fatal error in reading json file: %s", err))
	}

	var graphDecodeModel DecodeModel
	if err = json.Unmarshal(bytes, &graphDecodeModel); err != nil {
		panic(fmt.Errorf("Fatal error in Unmarshal(): %s", err))
	}

	return graphDecodeModel
}

func loadDistributions(graphDecodeModel DecodeModel, librarySize int) map[string]distribution.Distribution {
	dists := make(map[string]distribution.Distribution)
	for _, requestModel := range graphDecodeModel.RequestModels {
		params := make(map[string]float64)
		for index, key := range requestModel.ParameterKeys {
			params[key] = requestModel.Parameters[index]
		}
		switch requestModel.Model {
		case "gamma":
			dists[requestModel.ID] = gamma.New(params["K"], params["Theta"], librarySize)
		case "zipf":
			dists[requestModel.ID] = zipf.New(params["Skewness"], librarySize)
		}
	}
	return dists
}

func loadModel(graphDecodeModel DecodeModel) *Graph_t {
	graphDecodeModel.ModelRelations = make(map[interface{}]interface{})

	network := newGraph()
	network.model = graphDecodeModel
	network.loadOriginServer(graphDecodeModel)
	network.loadCacheServers(graphDecodeModel)
	network.loadLinks(graphDecodeModel)
	network.loadClients(graphDecodeModel)

	network.initRouter()

	if graphDecodeModel.Nodes[0].CacheAlgorithm == "iris" {
		network.initSpectrums()
	}

	return network
}

func LoadGraph(filename string) graph.Graph {
	graphDecodeModel := loadFile(filename)
	return loadModel(graphDecodeModel)
}
