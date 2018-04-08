package main

import "graph"
import "math/rand"
import "utils"
import "cache"
import "parser"
import "origin"
import "fmt"
import "cache/eviction/iris"
import "encoding/json"
import "cache/eviction/random"
import "cache/eviction/fifo"
import "cache/eviction/srrip"
import "cache/eviction/arc"
import "cache/eviction/lirs"
import "cache/eviction/lru"
import "cache/eviction/lruk"
import "cache/eviction/modifiedlru"
import "cache/eviction/admission"
import "cache/eviction/iclfu"
import "cache/eviction/wlfu"

type NodeReport map[string]interface{}
type LinkReport map[string]interface{}

type Graph_t struct {
	nodes           []*graph.Node_t
	links           []*graph.UnidirectionalLink_t
	clients         []*graph.Client_t
	router          graph.Router
	spectrumManager *SpectrumManager_t
	model           DecodeModel
}

func NewRouter(g *Graph_t) graph.Router {
	router := new(Router_t)
	router.graph = g
	return graph.Router(router)
}

func (g *Graph_t) SpectrumManager() graph.SpectrumManager {
	return g.spectrumManager
}

func (g *Graph_t) LibrarySize() int {
	librarySize := 0
	for _, origin := range g.originServers() {
		storageSize := origin.Entity().(*graph.ServerModel_t).Storage().Len()
		if librarySize < storageSize {
			librarySize = storageSize
		}
	}
	return librarySize
}

func (g *Graph_t) originServers() []*graph.Node_t {
	originServers := make([]*graph.Node_t, 0)
	for _, node := range g.nodes {
		if node.Entity().(*graph.ServerModel_t).IsOrigin() {
			originServers = append(originServers, node)
		}
	}
	return originServers
}

func (g *Graph_t) OriginServers() []graph.Node {
	originServers := make([]graph.Node, 0)
	for _, node := range g.nodes {
		if node.Entity().(*graph.ServerModel_t).IsOrigin() {
			originServers = append(originServers, graph.Node(node))
		}
	}
	return originServers
}

func (g *Graph_t) cacheServers() []*graph.Node_t {
	cacheServers := make([]*graph.Node_t, 0)
	for _, node := range g.nodes {
		if !node.Entity().(*graph.ServerModel_t).IsOrigin() {
			cacheServers = append(cacheServers, node)
		}
	}
	return cacheServers
}

func (g *Graph_t) Links() []graph.UnidirectionalLink {
	links := make([]graph.UnidirectionalLink, 0)
	for _, link := range g.links {
		links = append(links, graph.UnidirectionalLink(link))
	}
	return links
}

func (g *Graph_t) Clients() []graph.Client {
	clients := make([]graph.Client, 0)
	for _, c := range g.clients {
		clients = append(clients, graph.Client(c))
	}
	return clients
}

func newGraph() *Graph_t {
	g := new(Graph_t)
	g.nodes = make([]*graph.Node_t, 0)
	g.links = make([]*graph.UnidirectionalLink_t, 0)
	g.router = NewRouter(g)
	return g
}

func newNode(id string, entity interface{}) *graph.Node_t {
	node := new(graph.Node_t)
	node.SetID(id)
	node.SetEntity(entity)
	node.SetOutputLinks(make([]*graph.UnidirectionalLink_t, 0))
	node.SetInputLinks(make([]*graph.UnidirectionalLink_t, 0))
	return node
}

func (g *Graph_t) initSpectrums() {
}

func (g *Graph_t) initRouter() {
	g.router.Init()
}

func (g *Graph_t) addNode(node *graph.Node_t) {
	g.nodes = append(g.nodes, node)
}

func (g *Graph_t) detectNode(id string) *graph.Node_t {
	for _, node := range g.nodes {
		if node.ID() == id {
			return node
		}
	}
	return nil
}

func (g *Graph_t) detectLink(src, dst *graph.Node_t) *graph.UnidirectionalLink_t {
	for _, link := range g.links {
		if link.Src() == src && link.Dst() == dst {
			return link
		}
	}
	return nil
}

func (g *Graph_t) connect(src, dst *graph.Node_t, link *graph.UnidirectionalLink_t) {
	if g.detectLink(src, dst) != nil {
		return
	}
	link.SetSrc(src)
	link.SetDst(dst)
	src.SetOutputLinks_t(append(src.OutputLinks_t(), link))
	dst.SetInputLinks_t(append(dst.InputLinks_t(), link))
	g.links = append(g.links, link)
}

func (g *Graph_t) randomRequest() interface{} {
	return g.clients[rand.Intn(len(g.clients))].RandomRequest()
}

func (g *Graph_t) ResetCounters() {
	for _, node := range g.nodes {
		node.Entity().(*graph.ServerModel_t).Storage().ResetCount()
	}
	for _, link := range g.links {
		link.SetTraffic(0.0)
	}
}

func (g *Graph_t) MatchStoreData(best utils.BestSeparatorRanks, bitSize int) bool {
	if best.NetworkID != g.model.NetworkID {
		return false
	}
	if len(best.Ranks) != bitSize {
		return false
	}
	if best.RequestModel != g.model.Clients[0].RequestModelID {
		return false
	}
	for _, requestModel := range g.model.RequestModels {
		if requestModel.ID == best.RequestModel {
			for index := range best.RequestModelParams {
				if best.RequestModelParams[index] != requestModel.Parameters[index] {
					return false
				}
			}
		}
	}
	return true
}

func (g *Graph_t) GenerateBestSeparatorRanksData(bestSeparatorRanks []int) utils.BestSeparatorRanks {
	var best utils.BestSeparatorRanks
	best.NetworkID = g.model.NetworkID
	best.RequestModel = g.model.Clients[0].RequestModelID
	for _, requestModel := range g.model.RequestModels {
		if requestModel.ID == best.RequestModel {
			best.RequestModelParams = requestModel.Parameters
			break
		}
	}
	best.Ranks = bestSeparatorRanks
	return best
}

func (g *Graph_t) CacheAlgorithm() string {
	return g.model.Nodes[0].CacheAlgorithm
}

func (g *Graph_t) SpectrumCapacity() int {
	return 0
}

func newTrafficExpectation(network *Graph_t) *graph.TrafficExpectation_t {
	expectation := new(graph.TrafficExpectation_t)
	expectation.SetNetwork(network)
	expectation.Traffic = make(map[graph.UnidirectionalLink]float64)
	for _, link := range network.links {
		expectation.Traffic[link] = 0.0
	}
	return expectation
}

func (g *Graph_t) expectTraffic() *graph.TrafficExpectation_t {
	var forwardEntry graph.ForwardEntry
	g.fillCaches()
	expectation := newTrafficExpectation(g)

	// 各クライアントからすべてのコンテンツについてのリクエストを生成
	for _, client := range g.clients {
		for contentRank := 1; contentRank <= g.originServers()[0].Entity().(*graph.ServerModel_t).Storage().Len(); contentRank++ {
			fromID := client.Upstream().ID()
			// 1度の転送に発生する通信量を設定
			unitTraffic := client.Dist().PDF(contentRank)
			// リクエストを生成
			request := cache.ContentRequest{contentRank, make([]interface{}, 0), client.TrafficWeight()}
			// クライアント直上のサーバでヒットしたら通信量=0
			if client.Upstream().Storage().Exist(request.ContentKey) {
				continue
			}

			for {
				// 経由したサーバを更新
				loopDetected := false
				for _, nodeID := range request.XForwardedFor {
					if nodeID == fromID {
						loopDetected = true
					}
				}
				request.XForwardedFor = append(request.XForwardedFor, fromID)

				// 次の転送先を決定（キャッシュアルゴリズムによって経路制御方法を変える）
				if g.model.Nodes[0].CacheAlgorithm == "iris" && loopDetected == false && !parser.Options.UseShortestPath {
					forwardEntry = g.router.SelectForwardEntryBySpectrum(fromID, request)
				} else {
					forwardEntry = g.router.SelectForwardEntry(fromID, request)
				}

				// 転送リンクに通信量を加算
				expectation.Traffic[forwardEntry.Link()] += unitTraffic

				// 次のノードでヒットしたら経路制御終了
				if forwardEntry.Node().Entity().(*graph.ServerModel_t).Storage().Exist(request.ContentKey) {
					break
				}

				// リクエスト発生元を次のサーバにしてもう一度リクエストを転送
				fromID = forwardEntry.Node().ID()
			}
		}
	}

	return expectation
}

func (g *Graph_t) fillCaches() {
	switch g.model.Nodes[0].CacheAlgorithm {
	case "iris":
		for _, node := range g.cacheServers() {
			node.Entity().(*graph.ServerModel_t).Storage().(iris.Accessor).FillUp()
		}
	}
}

func (g *Graph_t) loadOriginServer(graphDecodeModel DecodeModel) {
	originServer := graph.NewServer(graphDecodeModel.Origin.ID, true)
	library := origin.NewLibrary(int(graphDecodeModel.Origin.LibrarySize))
	originServer.SetStorage(library)
	originNode := newNode(graphDecodeModel.Origin.ID, originServer)
	g.addNode(originNode)
	graphDecodeModel.ModelRelations[originNode] = graphDecodeModel.Origin
}

func (g *Graph_t) loadCacheServers(graphDecodeModel DecodeModel) {
	for _, nodeModel := range graphDecodeModel.Nodes {

		params := make(map[string]float64)
		for index, key := range nodeModel.ParameterKeys {
			params[key] = nodeModel.Parameters[index]
		}

		edgeServer := graph.NewServer(nodeModel.ID, false)

		var storage cache.Storage
		switch nodeModel.CacheAlgorithm {
		case "modifiedlru":
			storage = modifiedlru.New(int(params["Capacity"]), int(params["Jump"]))
		case "random":
			storage = random.New(int(params["Capacity"]))
		case "lru":
			storage = lru.New(int(params["Capacity"]))
		case "lruk":
			storage = lruk.New(int(params["Capacity"]), int(params["K"]))
		case "srrip":
			storage = srrip.New(int(params["Capacity"]), int(params["RRPVbit"]))
		case "lirs":
			storage = lirs.New(int(params["Capacity"]))
		case "arc":
			storage = arc.New(int(params["Capacity"]))
		case "fifo":
			storage = fifo.New(int(params["Capacity"]))
		case "iclfu":
			storage = iclfu.New(int(params["Capacity"]))
		case "windowlfu":
			storage = windowlfu.New(int(params["Capacity"]), int(params["Window"]))
		case "lfu":
			admissionList := make([]interface{}, 0)
			for rank := 0; rank < int(params["Capacity"]); rank++ {
				admissionList = append(admissionList, rank)
			}
			storage = admission.New(admissionList)
		case "iris":
			storage = NewIrisCache(int(params["Capacity"]), params["SpectrumRatio"])
		}
		edgeServer.SetStorage(storage)
		edgeServer.SetUpstreamRouter(g.router)

		cacheNode := newNode(nodeModel.ID, edgeServer)
		g.addNode(cacheNode)
		graphDecodeModel.ModelRelations[cacheNode] = nodeModel
	}
}

func (g *Graph_t) loadLinks(graphDecodeModel DecodeModel) {
	for _, linkModel := range graphDecodeModel.Links {
		SrcID := linkModel.EdgeNodeIds[0]
		DstID := linkModel.EdgeNodeIds[1]

		forwardLink := new(graph.UnidirectionalLink_t)
		forwardLink.SetCost(linkModel.Cost)
		g.connect(g.detectNode(SrcID), g.detectNode(DstID), forwardLink)
		graphDecodeModel.ModelRelations[forwardLink] = linkModel
		if linkModel.Bidirectional {
			inverseLink := new(graph.UnidirectionalLink_t)
			inverseLink.SetCost(linkModel.Cost)
			g.connect(g.detectNode(DstID), g.detectNode(SrcID), inverseLink)
			graphDecodeModel.ModelRelations[inverseLink] = linkModel
		}
	}
}

func (g *Graph_t) loadClients(graphDecodeModel DecodeModel) {
	librarySize := g.LibrarySize()
	dists := loadDistributions(graphDecodeModel, librarySize)
	for _, clientModel := range graphDecodeModel.Clients {
		cl := graph.NewClient(dists[clientModel.RequestModelID], g.detectNode(clientModel.UpstreamID).Entity().(*graph.ServerModel_t), clientModel.TrafficWeight)
		g.clients = append(g.clients, cl)
		graphDecodeModel.ModelRelations[cl] = clientModel
	}
}

func (g *Graph_t) clone() *Graph_t {
	return loadModel(g.model)
}

func (g *Graph_t) PlainReport() {
	utils.Print(fmt.Sprintln("Nodes:"))
	for index, node := range g.nodes {
		hit := node.Entity().(*graph.ServerModel_t).Storage().HitCount()
		miss := node.Entity().(*graph.ServerModel_t).Storage().MissCount()
		utils.Print(fmt.Sprintf("  [%3d]\t%s\t(access:%6d,\thit:%6d,\thit rate:%6.1f%%)\n", index, node.ID(), hit+miss, hit, float64(hit)/float64(hit+miss)*100.0))
	}
	utils.Print(fmt.Sprintln("Links:"))
	for index, link := range g.links {
		utils.Print(fmt.Sprintf("  [%3d]\t%s -> %s,\ttraffic:%6.1f\n", index, link.Src().ID(), link.Dst().ID(), link.Traffic()))
	}
}

func (g *Graph_t) JsonReport() {
	report := make(map[string]interface{})
	summary := make(map[string]interface{})
	originServerReport := make(map[string]interface{})
	cacheServerReports := make([]NodeReport, 0)
	linkReports := make([]LinkReport, 0)

	summary["internalTraffic"] = 0.0
	summary["originTraffic"] = 0.0

	origin := g.originServers()[0]
	originServerReport["id"] = origin.ID
	originServerReport["hit"] = origin.Entity().(*graph.ServerModel_t).Storage().HitCount()
	originServerReport["miss"] = origin.Entity().(*graph.ServerModel_t).Storage().MissCount()
	originServerReport["accesses"] = originServerReport["hit"].(int) + originServerReport["miss"].(int)
	originServerReport["hitrate"] = origin.Entity().(*graph.ServerModel_t).HitRate()
	originServerReport["caches"] = origin.Entity().(*graph.ServerModel_t).Storage().CacheList()

	for _, node := range g.cacheServers() {
		nodeReport := make(map[string]interface{})
		nodeReport["id"] = node.ID
		nodeReport["hit"] = node.Entity().(*graph.ServerModel_t).Storage().HitCount()
		nodeReport["miss"] = node.Entity().(*graph.ServerModel_t).Storage().MissCount()
		nodeReport["accesses"] = nodeReport["hit"].(int) + nodeReport["miss"].(int)
		nodeReport["hitrate"] = node.Entity().(*graph.ServerModel_t).HitRate()
		nodeReport["caches"] = node.Entity().(*graph.ServerModel_t).Storage().CacheList()
		cacheServerReports = append(cacheServerReports, nodeReport)
	}
	for _, link := range g.links {
		linkReport := make(map[string]interface{})
		linkReport["src"] = link.Src().ID()
		linkReport["dst"] = link.Dst().ID()
		linkReport["traffic"] = link.Traffic()

		if link.Src().ID() == originServerReport["id"] || link.Dst().ID() == originServerReport["id"] {
			summary["originTraffic"] = summary["originTraffic"].(float64) + link.Traffic()
		} else {
			summary["internalTraffic"] = summary["internalTraffic"].(float64) + link.Traffic()
		}

		linkReports = append(linkReports, linkReport)
	}
	summary["totalTraffic"] = summary["internalTraffic"].(float64) + summary["originTraffic"].(float64)

	report["OriginServer"] = originServerReport
	report["CacheServers"] = cacheServerReports
	report["Links"] = linkReports
	jsonString, _ := json.Marshal(report)
	utils.Print(fmt.Sprintln(string(jsonString)))
}
