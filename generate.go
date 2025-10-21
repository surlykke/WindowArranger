// Copyright (c) Christian Surlykke
//
// This file is part of the WindowArranger project.
// It is distributed under the GPL v2 license.
// Please refer to the GPL2 file for a copy of the license.
package main

import (
	_ "embed"
	"fmt"
	"io"
)

func output(out io.Writer, format string, args ...any) {
	fmt.Fprintf(out, format+"\n", args...)
}

func swaymsg(out io.Writer, criteria string, command string, args ...any) {
	if criteria != "" {
		criteria = "[" + criteria + "] "
	}
	var formatString = "swaymsg '" + criteria + command + "'"
	output(out, formatString, args...)
}

func doWait(out io.Writer, workspaces []Workspace, seconds uint) {
	if seconds == 0 {
		return
	}

	output(out, `
        # Takes a criteria as first argument and deadline as second
        function wait {
            while !  swaymsg "[$1] focus"> /dev/null; do 
                if [[ "$2" -lt "$(date +%%x)" ]]; then 
                    echo "Window specified by $1 did not appear before timeout"
                    exit 1
                fi
                sleep 0.1; 
            done
        }
    
        DEADLINE=$(( $(date +%%s) + %d ))
		`, seconds)

	for _, w := range workspaces {
		for _, criteria := range getAllCriteria(w.Children) {
			output(out, "wait '%s' $DEADLINE", criteria)
		}
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

func doNodeList(out io.Writer, nodes []*Node) {
	for _, node := range nodes {
		if node.Criteria == "" {
			var allC = getAllCriteria(node.Children)
			if len(allC) > 0 {
				swaymsg(out, allC[0], "focus")
				swaymsg(out, allC[0], "split v")
				swaymsg(out, allC[0], "layout %s", node.Layout)
				if len(allC) > 1 {
					swaymsg(out, allC[0], "mark current")
					for i := len(allC) - 1; i > 0; i-- {
						swaymsg(out, allC[i], "move to mark current")
					}
					swaymsg(out, "", "unmark current")
				}
			}
		}
	}
}

func translate(in io.Reader, out io.Writer, waitSeconds uint) {
	var workspaces = Parse(in)

	output(out, "#!/usr/bin/env bash")
	doWait(out, workspaces, waitSeconds)
	output(out, "# Move everything aside")
	swaymsg(out, "title=.*", "move workspace %d", len(workspaces)+1)

	for i, workspace := range workspaces {
		output(out, "# Workspace %d on %s", i+1, workspace.Output)
		var allC = getAllCriteria(workspace.Children)
		if len(allC) > 0 {
			swaymsg(out, allC[0], "move to workspace %d", i+1)
			swaymsg(out, allC[0], "focus")
			swaymsg(out, allC[0], "layout %s", workspace.Layout)
			if len(allC) > 1 {
				for j := len(allC) - 1; j > 0; j-- {
					swaymsg(out, allC[j], "move to workspace %d", i+1)
				}
			}
			swaymsg(out, "", "move workspace to output %s", workspace.Output)
		}
		doNodeList(out, workspace.Children)

	}
}
