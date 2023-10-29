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
		currentChar      rune = 0
		isInString       bool = false
		isInComment      bool = false
		runeCount        int  = 0
		currentLine      int  = 1
		currentLineStart int  = 0
	)

	var failIf = func(condition bool, msg string) {
		if condition {
			panic(fmt.Sprintf("%d,%d: %s", currentLine, runeCount-currentLineStart, msg))
		}
	}

	var isCurrentCharWhitespace = func() bool {
		return unicode.IsSpace(currentChar) || isInComment
	}

	var gotoNextChar = func() {
		r, w := utf8.DecodeRune(bytes)
		failIf(r == utf8.RuneError && w > 0, "Not valid utf-8")
		if isInComment {
			isInComment = r != '\n' && r != utf8.RuneError
		} else if isInString {
			failIf(r == utf8.RuneError, "String not terminated")
			if r == '\'' {
				isInString = false
			}
		} else {
			if r == '\'' {
				isInString = true
			} else if r == '#' {
				isInComment = true
			}
		}
		currentChar = r
		bytes = bytes[w:]	
	}

	var gotoNextNonWsChar = func() {
		for {
			gotoNextChar()
			if !isCurrentCharWhitespace() {
				break
			}
		}
	}

	var readString = func() string {

		var builder = strings.Builder{}

		for gotoNextChar(); isInString; gotoNextChar() {
			builder.WriteRune(currentChar)
		}
		gotoNextNonWsChar()
		return builder.String()
	}

	var readLayout func() *Node
	readLayout = func() *Node {
		var layoutMap = map[string]string{
			"H": "splith",
			"V": "splitv",
			"T": "tabbed",
			"S": "stacking",
		}
		var node = &Node{
			Layout: layoutMap[string(currentChar)],
		}
		failIf("" == node.Layout, "Layout type (H,V,T or S) expected")
		gotoNextNonWsChar()

		failIf(currentChar != '[', "'[' expected")
		gotoNextNonWsChar()

		for currentChar != ']' {
			if currentChar == '\'' {
				node.Children = append(node.Children, &Node{Criteria: readString()})
			} else {
				node.Children = append(node.Children, readLayout())
			}
		}
		gotoNextNonWsChar()
		return node
	}

	var readOutput = func() Workspace {
		var outputBuilder = strings.Builder{}

		failIf(!unicode.IsLetter(currentChar), "Output identifier expected")
		outputBuilder.WriteRune(currentChar)

		for gotoNextChar(); unicode.IsLetter(currentChar) || unicode.IsDigit(currentChar) || '-' == currentChar; gotoNextChar() {
			outputBuilder.WriteRune(currentChar)
		}
		if isCurrentCharWhitespace() {
			gotoNextNonWsChar()
		}
		failIf(':' != currentChar, "':' expected")
		gotoNextNonWsChar()
		return Workspace{Output: outputBuilder.String(), Node: readLayout()}
	}

	var workspaces []Workspace = nil
	gotoNextNonWsChar()
	for currentChar != utf8.RuneError {
		workspaces = append(workspaces, readOutput())
	}
	if len(workspaces) == 0 {
		panic("Empty configuration")
	}
	return workspaces
}
