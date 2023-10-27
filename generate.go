package main

import (
	_ "embed"
	"fmt"
)

var workspaceNo int = 1
var containerCount = 1

//go:embed functions.sh
var functions string

func Generate(workspaces []Workspace, waitSeconds uint) []string {
	const tempWorkspace = "window_arranger_temp_workspace"
	var program []string

	var add = func(format string, v ...interface{}) {
		program = append(program, fmt.Sprintf(format, v...))
	}

	var cmd = func(format string, v ...interface{}) {
		add("swaymsg '" + fmt.Sprintf(format, v...) + "'")
	}

	var createDummyWindow = func(node *Node) {
		var title = fmt.Sprintf("dummy_window_%02d", containerCount)
		containerCount = containerCount + 1
		node.Criteria = fmt.Sprintf("title=\"%s\"", title)
		add("dummywindow %s &", title)
	}

	var executeList func([]*Node)
	executeList = func(nodes []*Node) {
		for _, node := range nodes {
			if node.Children != nil {
				createDummyWindow(node)
				add("wait '%s'", node.Criteria)
			}
			cmd("[%s] move workspace %d; [%s] focus", node.Criteria, workspaceNo, node.Criteria)
		}
		for _, node := range nodes {
			if node.Children != nil {
				cmd("[%s] focus; splitv; layout %s", node.Criteria, node.Layout)
				executeList(node.Children)
			}
		}
	}

	add("#!/usr/bin/env bash")
	program = append(program, functions)
	if waitSeconds > 0 {
		add("waitWithDeadline $(( $(date +%%s) + %d )) %s", waitSeconds, collectCriteriaFromWorkspaces(workspaces))
	}
	cmd("[title=.*] move workspace %s", tempWorkspace)
	for _, workspace := range workspaces {
		executeList(workspace.Node.Children)
		cmd("[%s] focus", workspace.Node.Children[0].Criteria)
		cmd("focus parent")
		cmd("focus parent")
		cmd("layout %s", workspace.Node.Layout)
		cmd("move workspace to output %s", workspace.Output)
		workspaceNo = workspaceNo + 1
	}

	cmd(`[title="^dummy_window_"] kill`)
	add("if swaymsg '[workspace=%s] focus'; then", tempWorkspace)
	add("  swaymsg 'move workspace to output %s'", workspaces[0].Output)
	add("fi")
	
	return program
}

func collectCriteriaFromWorkspaces(workspaces []Workspace) string {
	var result = ""
	for _, workspace := range workspaces {
		result = result + collectCriteriaFromNodes(workspace.Node.Children)
	}
	return result
}

func collectCriteriaFromNodes(nodes []*Node) string {
	var result = ""
	for _, node := range nodes {
		if node.Children == nil {
			result = result + node.Criteria + " "
		} else {
			result = result + collectCriteriaFromNodes(node.Children)
		}
	}
	return result
}
