package graph

import (
	"github.com/mzohreva/GoGraphviz/graphviz"
	"log"
)

type IGraph interface {
	AddNode(label string) int
	AddEdge(from, to int, label string) int
	GetNodes() []Node
	GetEdges() []Edge
	GenerateImage(outputFileName string)
}

type Node struct {
	Label string
	Id    int
}

func NewNode(label string, id int) *Node {
	return &Node{Id: id, Label: label}
}

type Edge struct {
	From  int
	To    int
	Label string
	Id    int
}

func NewEdge(from, to int, label string, id int) *Edge {
	return &Edge{Label: label, From: from, To: to, Id: id}
}

type Graph struct {
	graph *graphviz.Graph
	nodes []Node
	edges []Edge
}

func NewGraph() *Graph {
	graph := &graphviz.Graph{}
	graph.MakeDirected()
	return &Graph{graph: graph, nodes: make([]Node, 0), edges: make([]Edge, 0)}
}

func (g *Graph) AddNode(label string) int {
	id := g.graph.AddNode(label)
	g.nodes = append(g.nodes, *NewNode(label, id))
	return id
}

func (g *Graph) AddEdge(from, to int, label string) int {
	id := g.graph.AddEdge(from, to, label)
	g.edges = append(g.edges, *NewEdge(from, to, label, id))
	return id
}

func (g *Graph) GetNodes() []Node {
	return g.nodes
}

func (g *Graph) GetEdges() []Edge {
	return g.edges
}

func (g *Graph) GenerateImage(outputFileName string) {
	err := g.graph.GenerateImage("dot", outputFileName+".png", "png")
	if err != nil {
		log.Fatal(err)
	}
}
