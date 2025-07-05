// Copyright (c) Christian Surlykke
//
// This file is part of the WindowArranger project.
// It is distributed under the GPL v2 license.
// Please refer to the GPL2 file for a copy of the license.
package main

import (
	_ "embed"
	"fmt"
)

const tempWorkspace = "window_arranger_temp_workspace"

var workspaceNo uint = 1
var containerCount uint = 1
var program []string

func add(format string, v ...any) {
	program = append(program, fmt.Sprintf(format, v...))
}

func cmd(format string, v ...any) {
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
			allCriteria = append(allCriteria, getAllCriteria(workspace.Children)...)
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
	var allCriteria = getAllCriteria(workspace.Children)
	if len(allCriteria) > 0 {
		cmd("[%s] move to workspace %d", allCriteria[0], workspaceNo)
		cmd("[%s] focus; layout %s", allCriteria[0], workspace.Layout)
		for i := len(allCriteria) - 1; i > 0; i-- {
			cmd("[%s] move to workspace %d", allCriteria[i], workspaceNo)
		}
	}
	for _, subNode := range workspace.Children {
		doSubNode(subNode)
	}
	add("")
}

func doSubNode(node *Node) {
	if node.Layout == "" {
		return
	}

	var allCriteria = getAllCriteria(node.Children)
	if len(allCriteria) > 0 {
		cmd(`[%s] focus; split v`, allCriteria[0])
		cmd(`[%s] layout %s; mark current`, allCriteria[0], node.Layout)
		for i := len(allCriteria) - 1; i > 0; i-- {
			cmd(`[%s] move to mark current`, allCriteria[i])
		}
		cmd(`unmark current`)
	}

	for _, subnode := range node.Children {
		doSubNode(subnode)
	}
}

func getAllCriteria(nodes []*Node) []string {
	var criteria = []string{}
	for _, n := range nodes {
		if n.Criteria != "" {
			criteria = append(criteria, n.Criteria)
		} else {
			criteria = append(criteria, getAllCriteria(n.Children)...)
		}
	}
	return criteria
}
