package main

import (
	_ "embed"
	"fmt"
)

var workspaceNo int = 1
var containerCount = 1

//go:embed scriptStart.sh
var scriptStart string

func Generate(workspaces []Workspace, waitSeconds uint) []string {
	const tempWorkspace = "window_arranger_temp_workspace"
	var dummywindowsCreated = false
	var program []string

	var add = func(format string, v ...interface{}) {
		program = append(program, fmt.Sprintf(format, v...))
	}

	var cmd = func(format string, v ...interface{}) {
		add("$MSG '" + fmt.Sprintf(format, v...) + "'")
	}

	var createDummyWindow = func(node *Node) {
		var title = fmt.Sprintf("dummy_window_%02d", containerCount)
		containerCount = containerCount + 1
		node.Criteria = fmt.Sprintf("title=\"%s\"", title)
		add("dummywindow %s &", title)
		dummywindowsCreated = true
	}

	var doNodeList func([]*Node)
	doNodeList = func(nodes []*Node) {
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
				doNodeList(node.Children)
			}
		}
	}

	program = []string{scriptStart}
	if waitSeconds > 0 {
		add("# Wait for all windows to be present")
		add("DEADLINE=$(( $(date +%%s) + %d ))", waitSeconds)
		add("wait  %s", collectCriteria(workspaces))
		add("DEADLINE=")
		add("")
	}

	add("# Move everything aside")
	cmd("[title=.*] move workspace %s", tempWorkspace)
	add("")

	for _, workspace := range workspaces {
		add("# Workspace %d on %s", workspaceNo, workspace.Output)
		doNodeList(workspace.Node.Children)
		cmd("[%s] focus", workspace.Node.Children[0].Criteria)
		cmd("focus parent")
		cmd("focus parent")
		cmd("layout %s", workspace.Node.Layout)
		cmd("move workspace to output %s", workspace.Output)
		workspaceNo = workspaceNo + 1
		add("")
	}

	if dummywindowsCreated {
		add("# Clean up")
		cmd(`[title="^dummy_window_"] kill`)
	}
	return program
}

func collectCriteria(workspaces []Workspace) string {
	var result = ""
	for _, workspace := range workspaces {
		result = result + collectCriteriaRecursively(workspace.Node.Children)
	}
	return result
}

func collectCriteriaRecursively(nodes []*Node) string {
	var result = ""
	for _, node := range nodes {
		if node.Children == nil {
			result = result + "'" + node.Criteria + "' "
		} else {
			result = result + collectCriteriaRecursively(node.Children)
		}
	}
	return result
}
