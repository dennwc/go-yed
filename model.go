package yed

import (
	"io"
	"strconv"
)

func NewFile(w io.Writer) *File {
	return &File{w: w, root: &Graph{}}
}

type File struct {
	w    io.Writer
	root *Graph

	lastEdge int
	edges    []*Edge
}

func (f *File) Graph() *Graph {
	return f.root
}

func (f *File) Edge(from, to *Node) *Edge {
	e := &Edge{f: f, id: "e" + strconv.Itoa(f.lastEdge), src: from, dst: to}
	f.lastEdge++
	f.edges = append(f.edges, e)
	return e
}

func (f *File) Close() error {
	w := newWriter(f.w)
	if err := w.WriteGraph(f.root); err != nil {
		return err
	}
	for _, e := range f.edges {
		if err := w.WriteEdge(e); err != nil {
			return err
		}
	}
	return w.Close()
}

type Graph struct {
	id string

	lastNode int
	sub      []*Node

	Description string
}

func (g *Graph) ID() string {
	return g.id
}

func (g *Graph) NewNode() *Node {
	id := "n" + strconv.Itoa(g.lastNode)
	g.lastNode++
	if g.id != "" {
		id = g.id + ":" + id
	}
	n := &Node{par: g, id: id}
	g.sub = append(g.sub, n)
	return n
}

type Node struct {
	id  string
	par *Graph
	sub *Graph

	Label       string
	Description string
	Style       *NodeStyle
}

func (n *Node) ID() string {
	return n.id
}

func (n *Node) SubGraph() *Graph {
	if n.sub == nil {
		n.sub = &Graph{id: n.id + ":"}
	}
	return n.sub
}

type Edge struct {
	f        *File
	id       string
	src, dst *Node

	Label       string
	Description string
	Style       *EdgeStyle
}

func (e *Edge) ID() string {
	return e.id
}

func (e *Edge) Source() *Node {
	return e.src
}

func (e *Edge) Target() *Node {
	return e.dst
}
