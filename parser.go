package main

import (
	"fmt"
)

type Workspace struct {
	Output string
	Node   *Node
}

type Node struct {
	Selector string
	Layout   string
	Children []*Node
}

func Parse(bytes []byte) []Workspace {
	var tokens = Tokenize(bytes) 
	var currentToken Token

	var next = func() {
		currentToken = <-tokens
	}

	var fail = func(message string) {
		panic(fmt.Sprintf("Fail at line %d, col %d: %s", currentToken.Line, currentToken.Col, message))
	}

	var readSingle = func(expected string) {
		if ok := currentToken.Text == expected; ok {
			next()
		} else {
			fail("'" + expected + "' expected")
		}
	}

	var readIdentifier = func() string {
		if currentToken.Type != Identifier {
			fail("Identifier expected")
		}

		var text = currentToken.Text
		next()
		return text
	}

	var readString = func() (string, bool) {
		if currentToken.Type == String {
			var text = currentToken.Text
			next()
			return text, true
		} else {
			return "", false
		}
	}

	var readLayout = func() string {
		var layout string
		switch currentToken.Text {
		case "H":
			layout = "splith"
		case "V":
			layout = "splitv"
		case "T":
			layout = "tabbed"
		case "S":
			layout = "stacking"
		default:
			fail("Layout type (H,V,T or S) expected")
		}
		next()
		return layout
	}


	var readNodeList func() *Node

	var readNode = func() *Node {
		if str, ok := readString(); ok {
			return &Node{Selector: str}
		} else {
			return readNodeList()
		}
	}

	readNodeList = func() *Node {
		var node = &Node{
			Layout: readLayout(),
		}
		readSingle("[")
		node.Children = append(node.Children, readNode()) // list may not be empty
		for currentToken.Text != "]" {
			node.Children = append(node.Children, readNode())
		}
		readSingle("]")
		return node
	}

	var readOutput = func() Workspace {
		var output = readIdentifier()
		readSingle(":")
		return Workspace{Output: output, Node: readNodeList()}
	}


	next()
	var workspaces []Workspace = nil
	for currentToken.Text != "" {
		workspaces = append(workspaces, readOutput())
	}
	return workspaces
}
