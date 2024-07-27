package main

import (
	_ "embed"
	"fmt"
)

const tempWorkspace = "window_arranger_temp_workspace"

//go:embed scriptStart.sh
var scriptStart string

var allCriteria string = ""
var workspaceNo uint = 1
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

func Generate(workspaces []Workspace, waitSeconds uint) []string {
	program = []string{scriptStart}
	if waitSeconds > 0 {
		wait(workspaces, waitSeconds)
	}

	add("# Move everything aside")
	cmd("[title=.*] move workspace %s", tempWorkspace)
	add("")

	for i, workspace := range workspaces {
		doWorkSpace(workspace, i+1)
	}
	return program
}

func doWorkSpace(workspace Workspace, workspaceNo int) {
	add("# Workspace %d on %s", workspaceNo, workspace.Output)
	cmd("workspace to output %d", workspaceNo)
	cmd("move workspace to output %s", workspace.Output)
	var leaves = getLeaves(workspace.Node)
	for _, leave := range getLeaves(workspace.Node) {
		cmd("[%s] move to workspace %d", leave, workspaceNo)
	}
	cmd("[%s] focus; layout %s", leaves[0], workspace.Node.Layout)

	for _, subNode := range workspace.Node.Children {
		doSubNode(subNode)
	}
	add("")
}

func doSubNode(node *Node) {
	if node.Layout == "" {
		return
	}
	var leaves = getLeaves(node)
	cmd(`[%s] focus; split v`, leaves[0])
	if node.Layout != "splitv" {
		cmd(`[%s] focus; layout %s`, leaves[0], node.Layout)
	}
	cmd(`[%s] mark current`, leaves[0])
	for i := len(leaves) - 1; i > 0; i-- {
		cmd(`[%s] move to mark current`, leaves[i])
	}
	cmd(`unmark current`)
	for _, subnode := range node.Children {
		doSubNode(subnode)
	}
}

func getLeaves(node *Node) []string {
	var leaves = make([]string, 0, 10)
	var walk func(*Node)
	walk = func(node *Node) {
		for _, child := range node.Children {
			if child.Criteria != "" {
				leaves = append(leaves, child.Criteria)
			} else {
				walk(child)
			}
		}
	}
	walk(node)
	return leaves
}
