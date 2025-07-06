package main

import (
	"fmt"
	"io"
)

type Parser struct {
	tokens   chan token
	curToken token
	ok       bool
}

func Parse(input io.Reader) []Workspace {
	var bytes, err = io.ReadAll(input)
	if err != nil {
		panic(err)
	}

	var tokens = make(chan token)

	go scan(tokens, []rune(string(bytes)))

	var curToken token
	var ok bool

	var getToken = func() {
		curToken, ok = <-tokens
	}

	var mustFind = func(token string) {
		getToken()
		if curToken.text != token {
			panic(fmt.Sprintf("Expected '%s' but found '%s'", token, curToken.text))
		}
	}

	var layoutMap = map[string]string{
		"H": "splith",
		"V": "splitv",
		"T": "tabbed",
		"S": "stacked"}

	var getLayout = func() string {
		if l, ok := layoutMap[curToken.text]; !ok {
			panic("Layout should be one of 'H', 'V', 'T' or 'S'")
		} else {
			return l
		}
	}

	var readNodes func() []*Node
	readNodes = func() []*Node {
		var nodes = make([]*Node, 0)
		for {
			getToken()
			if !ok {
				panic("Unexpected end of input")
			} else if curToken.text == "]" {
				break
			} else if curToken.Type == String {
				nodes = append(nodes, &Node{Criteria: curToken.text})
			} else {
				var node = &Node{
					Layout: getLayout(),
				}
				mustFind("[")
				node.Children = readNodes()
				nodes = append(nodes, node)
			}
		}

		return nodes
	}

	var readWorkSpace = func() Workspace {
		if curToken.Type != Id {
			panic("Any workspace definition should start with an output name")
		}
		var ws = Workspace{Output: curToken.text}
		mustFind(":")
		getToken()
		ws.Layout = getLayout()
		mustFind("[")
		ws.Children = readNodes()
		return ws
	}

	var workspaces = make([]Workspace, 0)
	for {
		getToken()
		if ok {
			workspaces = append(workspaces, readWorkSpace())
		} else {
			break
		}
	}
	return workspaces
}
