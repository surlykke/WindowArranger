package main

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

type Workspace struct {
	Output string
	Node   *Node
}

type Node struct {
	Criteria string
	Layout   string
	Children []*Node
}

func Parse(bytes []byte) []Workspace {
	var (
		runeCount        int = 0
		currentLine      int = 1
		currentLineStart int = 0
		layoutMap            = map[string]string{
			"H": "splith",
			"V": "splitv",
			"T": "tabbed",
			"S": "stacking",
		}
	)

	var failIf = func(condition bool, msg string) {
		if condition {
			panic(fmt.Sprintf("%d,%d: %s", currentLine, runeCount-currentLineStart, msg))
		}
	}

	var currentRune = func() rune {
		r, w := utf8.DecodeRune(bytes)
		failIf(r == utf8.RuneError && w > 0, "Not valid utf-8")
		return r
	}

	var skip = func() {
		_, w := utf8.DecodeRune(bytes)
		runeCount = runeCount + 1
		bytes = bytes[w:]
	}

	var skipWhitespace = func() {
		var inComment = false
		for {
			if unicode.IsSpace(currentRune()) || currentRune() == '#' || inComment {
				if currentRune() == '\n' {
					currentLine = currentLine + 1
					currentLineStart = runeCount
					inComment = false
				} else if currentRune() == '#' {
					inComment = true
				}
				skip()
			} else {
				break
			}
		}
	}

	var nextNonWsRune = func() rune {
		skipWhitespace()
		return currentRune()
	}

	var readString = func() string {
	
		var builder = strings.Builder{}
		
		for	skip(); '\'' != currentRune(); skip() { 
			failIf(currentRune() == utf8.RuneError , "String not terminated")
			builder.WriteRune(currentRune())
		}
		skip()
		return builder.String()
	}

	var readNodeList func() *Node
	readNodeList = func() *Node {
		var layout = layoutMap[string(nextNonWsRune())]
		failIf("" == layout, "Layout type (H,V,T or S) expected")
		skip()

		var node = &Node{
			Layout: layout,
		}

		failIf(nextNonWsRune() != '[', "'[' expected")
		skip()

		for nextNonWsRune() != ']' {
			if nextNonWsRune() == '\'' {
				node.Children = append(node.Children, &Node{Criteria: readString()})
			} else {
				node.Children = append(node.Children, readNodeList())
			}
		}
		skip()
		return node
	}

	var readOutput = func() Workspace {
		failIf(!unicode.IsLetter(nextNonWsRune()), "Output identifier expected")
		var outputBuilder = strings.Builder{}
		outputBuilder.WriteRune(currentRune())
		skip()

		for ; unicode.IsLetter(currentRune()) || unicode.IsDigit(currentRune()) || '-' == currentRune(); skip() {
			outputBuilder.WriteRune(currentRune())
		}

		failIf(':' != nextNonWsRune(), "':' expected")
		skip()
		return Workspace{Output: outputBuilder.String(), Node: readNodeList()}
	}

	var workspaces []Workspace = nil
	for nextNonWsRune() != utf8.RuneError {
		workspaces = append(workspaces, readOutput())
	}
	if len(workspaces) == 0 {
		panic("Empty configuration")
	}
	return workspaces
}

