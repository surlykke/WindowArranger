function wait {
    while (($#)); do 
        echo swaymsg "[$1] focus"
        while !  swaymsg "[$1] focus"> /dev/null; do 
            sleep 0.1; 
            if [[ -n "$DEADLINE" ]] && [[ "$DEADLINE" -lt $(date +%s) ]]; then 
                echo "Window specified by "$1" did not appear before timeout" >&2
                exit 1
            fi
        done
        shift
    done
}
