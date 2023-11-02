#!/usr/bin/env bash
export MSG=swaymsg
if [[ "i3" == "${XDG_SESSION_DESKTOP}" ]]; then 
	MSG=i3-msg
fi

function wait {
    while (($#)); do 
        while !  $MSG "[$1] focus"> /dev/null; do 
            sleep 0.1; 
            if [[ -n "$DEADLINE" ]] && [[ "$DEADLINE" -lt $(date +%s) ]]; then 
                echo "Window specified by "$1" did not appear before timeout" >&2
                exit 1
            fi
        done
        shift
    done
}

