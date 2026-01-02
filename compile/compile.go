// Copyright (c) Christian Surlykke
//
// This file is part of the WindowArranger project.
// It is distributed under the GPL v2 license.
// Please refer to the GPL2 file for a copy of the license.
package compile

import (
	"fmt"
	"io"
	"slices"
	"unicode"

	"github.com/goccy/go-yaml"
	"github.com/surlykke/WindowArranger/sway"
)

// Yaml is read into these two
type Layout struct {
	Monitors     []Monitor
	Postcommands []string
}

type Monitor struct {
	Name       string
	Make       string
	Model      string
	Serial     string
	Workspaces []string
	Posx       *int
	Posy       *int
	Scale      *float32
}

// workspacespecs (eg.  H[V['title=T1' 'title=T2'] H['title=T3' 'title=T4']]) are parsed to a treestructure
type Node struct {
	layout      string
	subnodes    []Node
	allcriteria []string // All criteria found within this node
}

func CompileConfig(input io.Reader) (program []string, criteria []string) {

	var add = func(line string, args ...any) {
		program = append(program, fmt.Sprintf(line, args...))
	}

	var layout Layout
	if confBytes, err := io.ReadAll(input); err != nil {
		panic(err)
	} else if err := yaml.Unmarshal(confBytes, &layout); err != nil {
		panic(err)
	} else {
		var numWorkspaces = 0
		var outputs = sway.GetOutputs()
		for i, m := range layout.Monitors {
			numWorkspaces += len(m.Workspaces)
			if m.Name == "" {
				layout.Monitors[i].Name = adjustName(m, outputs)
			}
		}
		add("[title=.*] move to workspace %d", numWorkspaces+1)

		var workspaceNo = 1
		for _, monitorLayout := range layout.Monitors {
			for _, workspaceSpec := range monitorLayout.Workspaces {
				var node = parse(workspaceSpec)
				generate(node, workspaceNo, true, add)
				add("[%s] focus", node.allcriteria[0])
				add("move workspace to %s", monitorLayout.Name)
				criteria = node.allcriteria
				workspaceNo++
			}

			if monitorLayout.Posx != nil && monitorLayout.Posy != nil {
				add("output %s position %d %d", monitorLayout.Name, *monitorLayout.Posx, *monitorLayout.Posy)
			} else if monitorLayout.Posx != nil || monitorLayout.Posy != nil {
				panic(monitorLayout.Name + ": When giving position, both posx and posy must be defined")
			}

			if monitorLayout.Scale != nil {
				add("output %s scale %.1f", monitorLayout.Name, *monitorLayout.Scale)
			}

		}

		for _, postCommand := range layout.Postcommands {
			add(postCommand)
		}
		return
	}
}

func adjustName(m Monitor, outputs []sway.Output) string {
	if m.Name != "" {
		return m.Name
	} else if m.Make == "" && m.Model == "" && m.Serial == "" {
		panic("Neither name, make, model or serial set for monitor")
	} else {
		for _, o := range outputs {
			if (m.Make == "" || m.Make == o.Make) &&
				(m.Model == "" || m.Model == o.Model) &&
				(m.Serial == "" || m.Serial == o.Serial) {
				return o.Name
			}
		}
	}
	panic(fmt.Sprintf("No match for monitor, make: %s, model: %s, serial %s", m.Make, m.Model, m.Serial))
}

const (
	H uint8 = iota
	V
	T
	S
	Lb
	Rb
	Str
	END
)

func parse(workspaceSpec string) Node {
	var runes, pos = append([]rune(workspaceSpec), 0), -1

	var curToken uint8
	var curText string

	var readString = func() string {
		var start = pos
		for pos++; runes[pos] != 0; pos++ {
			if runes[pos] == '\'' {
				return string(runes[start+1 : pos])
			}
		}
		panic("Runaway string")
	}

	var nextToken = func() uint8 {
		for pos++; unicode.IsSpace(runes[pos]); pos++ {
		}

		switch runes[pos] {
		case 'H':
			curToken, curText = H, "splith"
		case 'V':
			curToken, curText = V, "splitv"
		case 'T':
			curToken, curText = T, "tabbed"
		case 'S':
			curToken, curText = S, "stacked"
		case '[':
			curToken, curText = Lb, ""
		case ']':
			curToken, curText = Rb, ""
		case '\'':
			curToken, curText = Str, readString()
		case 0:
			curToken, curText = END, ""
		default:
			panic("Unexpected character: " + string(runes[pos]))
		}

		return curToken
	}

	var readNode func() Node

	readNode = func() Node {
		if !slices.Contains([]uint8{H, V, T, S}, curToken) {
			panic(fmt.Sprintf("H, V, T or S expected at: '%s':%d", string(runes), pos))
		}

		var node = Node{layout: curText}
		if nextToken() != Lb {
			panic(fmt.Sprintf("[ expected, got: %d\n", curToken))
		}
		for {
			switch nextToken() {
			case Rb:
				if len(node.allcriteria) == 0 {
					panic("empty node")
				}
				return node
			case Str:
				node.allcriteria = append(node.allcriteria, curText)
			default:
				node.subnodes = append(node.subnodes, readNode())
				node.allcriteria = append(node.allcriteria, node.subnodes[len(node.subnodes)-1].allcriteria...)
			}
		}
	}

	nextToken()
	var node = readNode()
	if nextToken() != END {
		panic("Workspace '" + workspaceSpec + "': trailing characters: " + string(runes[pos:]))
	}

	return node
}

func generate(node Node, workspaceNo int, root bool, add func(string, ...any)) {
	if root {
		add("[%s] focus, move to workspace %d, layout %s", node.allcriteria[0], workspaceNo, node.layout)
	} else {
		add("[%s] focus, split v, layout %s", node.allcriteria[0], node.layout)
	}

	if len(node.allcriteria) > 1 {
		add("[%s] focus, mark current", node.allcriteria[0])
		for i := len(node.allcriteria) - 1; i > 0; i-- {
			add("[%s] move to mark current", node.allcriteria[i])
		}
		add("unmark current")
	}

	for _, subnode := range node.subnodes {
		generate(subnode, workspaceNo, false, add)
	}
}
