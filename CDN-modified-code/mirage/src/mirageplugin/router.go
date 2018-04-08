package main

import "fmt"
import "parser"
import "cache"
import "utils"
import "math"
import "graph"

type Router_t struct {
	graph        *Graph_t
	routingTable map[string]map[string][]graph.ForwardEntry
}

func newSpectrumManager(bitSize int, network *Graph_t) *SpectrumManager_t {
	manager := new(SpectrumManager_t)
	return manager
}

func (router *Router_t) RoutingTable() map[string]map[string][]graph.ForwardEntry {
	return router.routingTable
}

func (router *Router_t) dijkstra(dstNode *graph.Node_t) {
	router.resetDijkstraVariables()
	dstNode.SetDijkstraCost(0.0)
	var doneNode *graph.Node_t

	for {
		doneNode = nil
		for _, node := range router.graph.nodes {
			if node.DijkstraDone() || node.DijkstraCost() < 0 {
				continue
			}
			if doneNode == nil || node.DijkstraCost() < doneNode.DijkstraCost() {
				doneNode = node
			}
		}
		if doneNode == nil {
			break
		}

		doneNode.SetDijkstraDone(true)
		for i := 0; i < len(doneNode.OutputLinks()); i++ {
			to := doneNode.OutputLinks()[i].Dst()
			cost := doneNode.DijkstraCost() + doneNode.OutputLinks()[i].Cost()
			if to.DijkstraCost() < 0 || cost < to.DijkstraCost() {
				to.SetDijkstraCost(cost)
			}
		}

	}
}

func (router *Router_t) createRoutingTable() {
	router.routingTable = make(map[string]map[string][]graph.ForwardEntry)

	for _, dstNode := range router.graph.nodes {
		router.dijkstra(dstNode)

		for _, srcNode := range router.graph.nodes {
			if router.routingTable[srcNode.ID()] == nil {
				router.routingTable[srcNode.ID()] = make(map[string][]graph.ForwardEntry)
			}
			if router.routingTable[srcNode.ID()][dstNode.ID()] == nil {
				router.routingTable[srcNode.ID()][dstNode.ID()] = make([]graph.ForwardEntry, 0)
			}
			if srcNode == dstNode {
				continue
			}
			minCost := math.Inf(0)

			for _, link := range srcNode.InputLinks() {
				if minCost > srcNode.Opposite(link).DijkstraCost() {
					minCost = srcNode.Opposite(link).DijkstraCost()
				}
			}
			for _, link := range srcNode.InputLinks() {
				if srcNode.Opposite(link).DijkstraCost() == minCost {
					forwardEntry := new(graph.ForwardEntry_t)
					forwardEntry.SetNode(srcNode.Opposite(link))
					forwardEntry.SetLink(link)
					forwardEntry.SetCost(srcNode.DijkstraCost())
					router.routingTable[srcNode.ID()][dstNode.ID()] =
						append(router.routingTable[srcNode.ID()][dstNode.ID()], forwardEntry)
				}
			}
		}

	}
}

func (router *Router_t) resetDijkstraVariables() {
	for _, node := range router.graph.nodes {
		node.SetDijkstraDone(false)
		node.SetDijkstraCost(-1.0)
	}
}

func (router *Router_t) inspectForwardEntry(entry graph.ForwardEntry) {
	utils.DebugPrint(fmt.Sprintf("    <next: %s, via_link: %s -> %s, cost:%f>\n", entry.Node().ID(), entry.Link().Src().ID(), entry.Link().Dst().ID(), entry.Cost()))
}

func (router *Router_t) inspectRoutingTable(fromNode *graph.Node_t) {
	utils.DebugPrint(fmt.Sprintf("From: %s\n", fromNode.ID()))
	for _, destNode := range router.graph.nodes {
		utils.DebugPrint(fmt.Sprintf("  Dest: %s\n", destNode.ID()))
		for _, forwardEntry := range router.routingTable[fromNode.ID()][destNode.ID()] {
			router.inspectForwardEntry(forwardEntry)
		}
	}
}

func (router *Router_t) inspectAllRoutingTables() {
	for _, fromNode := range router.graph.nodes {
		router.inspectRoutingTable(fromNode)
	}
}

func (router *Router_t) Init() {
	router.createRoutingTable()
}

func (router *Router_t) SelectForwardEntryBySpectrum(fromID string, request cache.ContentRequest) graph.ForwardEntry {
	return router.SelectForwardEntry(fromID, request)
}

func (router *Router_t) SelectForwardEntry(fromID string, request cache.ContentRequest) graph.ForwardEntry {
	srcNode := router.graph.detectNode(fromID)
	originServers := router.graph.originServers()
	forwardEntries := make([]graph.ForwardEntry, 0)

	for _, origin := range originServers {
		for _, forwardEntry := range router.routingTable[srcNode.ID()][origin.ID()] {
			forwardEntries = append(forwardEntries, forwardEntry)
		}
	}

	minCost := math.Inf(0)
	for _, forwardEntry := range forwardEntries {
		if minCost > forwardEntry.Cost() {
			minCost = forwardEntry.Cost()
		}
	}

	minCostEntries := make([]graph.ForwardEntry, 0)
	for _, forwardEntry := range forwardEntries {
		if forwardEntry.Cost() == minCost {
			minCostEntries = append(minCostEntries, forwardEntry)
		}
	}

	if len(minCostEntries) > 0 {
		return minCostEntries[0]
	}
	return nil
}

func (router *Router_t) ForwardRequest(fromID string, request cache.ContentRequest) interface{} {
	var forwardEntry graph.ForwardEntry
	if router.graph.model.Nodes[0].CacheAlgorithm == "iris" && !parser.Options.UseShortestPath {
		forwardEntry = router.SelectForwardEntryBySpectrum(fromID, request)
	} else {
		forwardEntry = router.SelectForwardEntry(fromID, request)
	}
	surrogate := forwardEntry.Node().Entity().(*graph.ServerModel_t)
	forwardEntry.Link().(*graph.UnidirectionalLink_t).AddTraffic(request.TrafficWeight)
	return surrogate.AcceptRequest(request)
}
