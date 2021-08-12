package main

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"github.com/GaudiestTooth17/irn-sim/network"
	"gonum.org/v1/gonum/graph"
)

// Read a file in GraphML
func readFile(filename string) *network.AdjacencyList {
	fileContents, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	tokens := tokenizeGMLString(string(fileContents))
	nodes, edges := parseGraph(tokens)
	sort.Slice(nodes, func(i, j int) bool { return nodes[i].ID() < nodes[j].ID() })

	adjList := make(map[int64][]graph.Node)
	for _, u := range nodes {
		adjList[u.ID()] = make([]graph.Node, 0)
	}
	for _, e := range edges {
		u := e.From()
		v := e.To()
		adjList[u.ID()] = append(adjList[u.ID()], v)
		adjList[v.ID()] = append(adjList[v.ID()], u)
	}
	net := network.NewAdjacencyList(nodes, adjList)
	return net
}

func tokenizeGMLString(gml string) []token {
	tokens := make([]token, 0)
	var sBuilder strings.Builder
	lineNum := 1
	columnNum := 0
	for _, c := range gml {
		columnNum++
		if !unicode.IsSpace(c) {
			sBuilder.WriteRune(c)
		} else if sBuilder.Len() > 0 {
			content := sBuilder.String()
			newToken := token{content, lineNum, columnNum - len(content)}
			tokens = append(tokens, newToken)
			sBuilder.Reset()
			if c == '\n' {
				lineNum++
				columnNum = 0
			}
		}
	}
	return tokens
}

// return the nodes and the edges
func parseGraph(tokens []token) ([]graph.Node, []graph.Edge) {
	tokens = match(tokens, "graph")
	tokens = match(tokens, "[")
	nodes, tokens := parseNodeList(tokens)
	edges, tokens := parseEdgeList(tokens)
	match(tokens, "]")
	return nodes, edges
}

func parseNodeList(tokens []token) ([]graph.Node, []token) {
	nodes := make([]graph.Node, 0)
	var u graph.Node
	for tokens[0].content != "edge" {
		u, tokens = parseNode(tokens)
		nodes = append(nodes, u)
	}
	return nodes, tokens
}

func parseNode(tokens []token) (graph.Node, []token) {
	tokens = match(tokens, "node")
	tokens = match(tokens, "[")
	tokens = match(tokens, "id")
	u, err := strconv.Atoi(tokens[0].content)
	tokens = tokens[1:]
	if err != nil {
		panic(err)
	}
	// skip the rest of the data
	for tokens[0].content != "]" {
		tokens = tokens[1:]
	}
	tokens = match(tokens, "]")
	return network.NewVertex(int64(u)), tokens
}

func parseEdgeList(tokens []token) ([]graph.Edge, []token) {
	edges := make([]graph.Edge, 0)
	var e graph.Edge
	for tokens[0].content != "]" {
		e, tokens = parseEdge(tokens)
		edges = append(edges, e)
	}
	return edges, tokens
}

func parseEdge(tokens []token) (graph.Edge, []token) {
	tokens = match(tokens, "edge")
	tokens = match(tokens, "[")

	tokens = match(tokens, "source")
	u, err := strconv.Atoi(tokens[0].content)
	tokens = tokens[1:]
	if err != nil {
		panic(err)
	}

	tokens = match(tokens, "target")
	v, err := strconv.Atoi(tokens[0].content)
	tokens = tokens[1:]
	if err != nil {
		panic(err)
	}

	// skip the rest of the data
	for tokens[0].content != "]" {
		tokens = tokens[1:]
	}
	tokens = match(tokens, "]")

	e := network.NewLink(network.NewVertex(int64(u)), network.NewVertex(int64(v)))
	return e, tokens
}

func match(tokens []token, expected string) []token {
	if tokens[0].content != expected {
		bad := tokens[0]
		panic(fmt.Errorf("parsing error: Expected '%s' Got '%s' at line: %d and col: %d",
			expected, bad.content, bad.line, bad.column))
	}
	return tokens[1:]
}

type token struct {
	content string
	line    int
	column  int
}
