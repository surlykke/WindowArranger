package main

import (
	_ "embed"
	"fmt"
)

const tempWorkspace = "window_arranger_temp_workspace"

//go:embed scriptStart.sh
var scriptStart string

var allCriteria    string = ""
var workspaceNo    uint = 1
var containerCount uint = 1
var program []string

func add(format string, v ...interface{}) {
	program = append(program, fmt.Sprintf(format, v...))
}

func cmd(format string, v ...interface{}) {
	add("$MSG '" + fmt.Sprintf(format, v...) + "'")
}

func collectCriteriaRecursively(node *Node) {
	if node.Children == nil {
		allCriteria = allCriteria + "'" + node.Criteria + "' "
	} else {
		for _, child := range node.Children {
			collectCriteriaRecursively(child)
		}
	}
}

func wait(workspaces []Workspace, waitSeconds uint) {
	for _, workspace := range workspaces {
		collectCriteriaRecursively(workspace.Node)
	}

	add("# Wait for all windows to be present")
	add("DEADLINE=$(( $(date +%%s) + %d ))", waitSeconds)
	add("wait  %s", allCriteria)
	add("DEADLINE=")
	add("")
}

func createDummyWindow() string {
	var title = fmt.Sprintf("dummy_window_%02d", containerCount)
	containerCount = containerCount + 1
	var criteria = fmt.Sprintf("title=%s", title)
	add("dummywindow %s &", title)
	add("wait '%s'", criteria)
	cmd("[%s] move workspace %s", criteria, tempWorkspace)
	return criteria
}

func createDummyWindows(node *Node) {
	if node.Children != nil {
		node.Criteria = createDummyWindow()
		for _, child := range node.Children {
			createDummyWindows(child)
		}
	}
}

func doNode(node *Node) {
	cmd("[%s] focus", node.Criteria)
	for _, child := range node.Children {
		cmd("[%s] move workspace %d", child.Criteria, workspaceNo)
		cmd("[%s] focus", child.Criteria)
	}

	for _, child := range node.Children {
		if child.Children != nil {
			cmd("[%s] focus; splitv; layout %s", child.Criteria, child.Layout)
			doNode(child)
		}
	}
}

func Generate(workspaces []Workspace, waitSeconds uint) []string {
	program = []string{scriptStart}
	if waitSeconds > 0 {
		wait(workspaces, waitSeconds)
	}

	add("# Move everything aside")
	cmd("[title=.*] move workspace %s", tempWorkspace)
	add("")

	for _, workspace := range workspaces {
		add("# Workspace %d on %s", workspaceNo, workspace.Output)
		createDummyWindows(workspace.Node)
		cmd("[%s] move workspace %d", workspace.Node.Criteria, workspaceNo)
		cmd("[%s] focus; focus parent; layout %s", workspace.Node.Criteria, workspace.Node.Layout)
		cmd("move workspace to output %s", workspace.Output)
		doNode(workspace.Node)
		cmd(`[title="^dummy_window_"] kill`)
		add("")
		workspaceNo = workspaceNo + 1
	}

	return program
}
