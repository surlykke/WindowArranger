package main

import (
	_ "embed"
	"fmt"
)

var workspaceNo int = 1
var containerCount = 1

//go:embed wait.sh
var waitFunction string

func Generate(workspaces []Workspace) []string {
	var program []string

	var add = func(line string) {
		program = append(program, line)
	}

	var cmd = func(format string, v ...interface{}) {
		add("swaymsg '" + fmt.Sprintf(format, v...) + "'")
	}

	var createDummyWindow = func(node *Node) {
		var title = fmt.Sprintf("dummy_window_%02d", containerCount)
		containerCount = containerCount + 1
		node.Selector = fmt.Sprintf("title=\"%s\"", title)
		add(fmt.Sprintf("dummywindow %s &", title))
	}

	
	var executeList func(nodes []*Node)
	executeList = func(nodes []*Node) {
		for _, node := range nodes {
			if node.Children != nil {
				createDummyWindow(node)
				add(fmt.Sprintf("wait '%s'", node.Selector))
			}
			cmd("[%s] move workspace %d; [%s] focus", node.Selector, workspaceNo, node.Selector)
		}
		for _, node := range nodes {
			if node.Children != nil {
				cmd("[%s] focus; splitv; layout %s", node.Selector, node.Layout)
				executeList(node.Children)
			}
		}
	}

	add("#!/usr/bin/env bash")
	add(waitFunction)
	cmd("[title=.*] move workspace temp")
	for _, workspace := range workspaces {
		executeList(workspace.Node.Children)
		cmd("[%s] focus", workspace.Node.Children[0].Selector)
		cmd("focus parent")
		cmd("focus parent")
		cmd("layout %s", workspace.Node.Layout)
		cmd("move workspace to output %s", workspace.Output)
		workspaceNo = workspaceNo + 1
	}

	cmd(`[title="^dummy_window_"] kill`)
	return program
}
