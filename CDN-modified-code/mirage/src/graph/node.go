package graph

type Node_t struct {
	id           string
	entity       interface{}
	outputLinks  []*UnidirectionalLink_t
	inputLinks   []*UnidirectionalLink_t
	dijkstraDone bool
	dijkstraCost float64
}

func (node *Node_t) SetDijkstraCost(cost float64) {
	node.dijkstraCost = cost
}

func (node *Node_t) DijkstraCost() float64 {
	return node.dijkstraCost
}

func (node *Node_t) DijkstraDone() bool {
	return node.dijkstraDone
}

func (node *Node_t) SetDijkstraDone(done bool) {
	node.dijkstraDone = done
}

func (node *Node_t) SetID(id string) {
	node.id = id
}

func (node *Node_t) SetEntity(entity interface{}) {
	node.entity = entity
}

func (node *Node_t) SetOutputLinks(outputLinks []*UnidirectionalLink_t) {
	node.outputLinks = outputLinks
}

func (node *Node_t) SetInputLinks(inputLinks []*UnidirectionalLink_t) {
	node.inputLinks = inputLinks
}

func (node *Node_t) ID() string {
	return node.id
}

func (node *Node_t) Entity() interface{} {
	return node.entity
}

func (node *Node_t) OutputLinks_t() []*UnidirectionalLink_t {
	return node.outputLinks
}

func (node *Node_t) InputLinks_t() []*UnidirectionalLink_t {
	return node.inputLinks
}

func (node *Node_t) SetOutputLinks_t(outputLinks []*UnidirectionalLink_t) {
	node.outputLinks = outputLinks
}

func (node *Node_t) SetInputLinks_t(inputLinks []*UnidirectionalLink_t) {
	node.inputLinks = inputLinks
}

func (node *Node_t) OutputLinks() []UnidirectionalLink {
	links := make([]UnidirectionalLink, 0)
	for _, link := range node.outputLinks {
		links = append(links, UnidirectionalLink(link))
	}
	return links
}

func (node *Node_t) InputLinks() []UnidirectionalLink {
	links := make([]UnidirectionalLink, 0)
	for _, link := range node.inputLinks {
		links = append(links, UnidirectionalLink(link))
	}
	return links
}

func (node *Node_t) Opposite(link UnidirectionalLink) *Node_t {
	if link.(*UnidirectionalLink_t).src == node {
		return link.Dst().(*Node_t)
	} else if link.Dst().(*Node_t) == node {
		return link.Src().(*Node_t)
	}
	return nil
}

func (node *Node_t) InputAdjacent() []*Node_t {
	adjacent := make([]*Node_t, 0)
	for _, link := range node.inputLinks {
		if link.dst == node {
			adjacent = append(adjacent, link.src)
		}
	}
	return adjacent
}

func (node *Node_t) outputAdjacent() []*Node_t {
	adjacent := make([]*Node_t, 0)
	for _, link := range node.outputLinks {
		if link.src == node {
			adjacent = append(adjacent, link.dst)
		}
	}
	return adjacent
}
