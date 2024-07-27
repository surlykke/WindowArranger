package main

import (
	_ "embed"
	"fmt"
)

const tempWorkspace = "window_arranger_temp_workspace"

var workspaceNo uint = 1
var containerCount uint = 1
var program []string

func add(format string, v ...interface{}) {
	program = append(program, fmt.Sprintf(format, v...))
}

func cmd(format string, v ...interface{}) {
	add("swaymsg '" + fmt.Sprintf(format, v...) + "'")
}

func wait(allCriteria []string, waitSeconds uint) {
	add(`# Wait for all windows to be present`)

	add(`# Takes a criteria as first argument and deadline as second`)
	add(`function wait {`)
	add(`	while !  swaymsg "[$1] focus"> /dev/null; do `)
	add(`		if [[ "$2" -lt "$(date +%%x)" ]]; then `)
	add(`			echo "Window specified by $1 did not appear before timeout"`)
	add(`			exit 1`)
	add(`		fi`)
	add(`		sleep 0.1; `)
	add(`	done`)
	add(`}`)

	add("DEADLINE=$(( $(date +%%s) + %d ))", waitSeconds)
	for _, criteria := range allCriteria {
		add("wait '%s' $DEADLINE", criteria)
	}
	add("")
}

func Generate(workspaces []Workspace, waitSeconds uint) []string {
	program = []string{`#!/usr/bin/env bash`}
	if waitSeconds > 0 {
		var allCriteria = make([]string, 0, 20)
		for _, workspace := range workspaces {
			allCriteria = append(allCriteria, getAllCriteria(workspace.Node)...)
		}
		wait(allCriteria, waitSeconds)
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
	var allCriteria = getAllCriteria(workspace.Node)
	for _, criteria := range getAllCriteria(workspace.Node) {
		cmd("[%s] move to workspace %d", criteria, workspaceNo)
	}
	cmd("[%s] focus; layout %s", allCriteria[0], workspace.Node.Layout)

	for _, subNode := range workspace.Node.Children {
		doSubNode(subNode)
	}
	add("")
}

func doSubNode(node *Node) {
	if node.Layout == "" {
		return
	}
	var allCriteriaForNode = getAllCriteria(node)
	cmd(`[%s] focus; split v`, allCriteriaForNode[0])
	if node.Layout != "splitv" {
		cmd(`[%s] focus; layout %s`, allCriteriaForNode[0], node.Layout)
	}
	cmd(`[%s] mark current`, allCriteriaForNode[0])
	for i := len(allCriteriaForNode) - 1; i > 0; i-- {
		cmd(`[%s] move to mark current`, allCriteriaForNode[i])
	}
	cmd(`unmark current`)
	for _, subnode := range node.Children {
		doSubNode(subnode)
	}
}

func getAllCriteria(node *Node) []string {
	var allCriteria = make([]string, 0, 10)
	var walk func(*Node)
	walk = func(node *Node) {
		for _, child := range node.Children {
			if child.Criteria != "" {
				allCriteria = append(allCriteria, child.Criteria)
			} else {
				walk(child)
			}
		}
	}
	walk(node)
	return allCriteria
}
