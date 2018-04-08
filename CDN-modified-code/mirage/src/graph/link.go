package graph

type UnidirectionalLink_t struct {
	src     *Node_t
	dst     *Node_t
	cost    float64
	traffic float64
}

func (link *UnidirectionalLink_t) SetCost(cost float64) {
	link.cost = cost
}

func (link *UnidirectionalLink_t) Cost() float64 {
	return link.cost
}

func (link *UnidirectionalLink_t) SetSrc(src *Node_t) {
	link.src = src
}

func (link *UnidirectionalLink_t) SetDst(dst *Node_t) {
	link.dst = dst
}

func (link *UnidirectionalLink_t) SetTraffic(traffic float64) {
	link.traffic = traffic
}

func (link *UnidirectionalLink_t) AddTraffic(traffic float64) {
	link.traffic += traffic
}

func (link *UnidirectionalLink_t) Traffic() float64 {
	return link.traffic
}

func (link *UnidirectionalLink_t) Src() Node {
	return link.src
}

func (link *UnidirectionalLink_t) Dst() Node {
	return link.dst
}
