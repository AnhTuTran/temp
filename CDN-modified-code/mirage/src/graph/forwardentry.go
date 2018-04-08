package graph

type ForwardEntry_t struct {
	node Node
	link UnidirectionalLink
	cost float64
}

type ForwardEntry interface {
	Node() Node
	Link() UnidirectionalLink
	Cost() float64
}

func (entry *ForwardEntry_t) SetNode(node Node) {
	entry.node = node
}

func (entry *ForwardEntry_t) SetLink(link UnidirectionalLink) {
	entry.link = link
}

func (entry *ForwardEntry_t) SetCost(cost float64) {
	entry.cost = cost
}

func (entry *ForwardEntry_t) Node() Node {
	return entry.node
}

func (entry *ForwardEntry_t) Link() UnidirectionalLink {
	return entry.link
}

func (entry *ForwardEntry_t) Cost() float64 {
	return entry.cost
}
