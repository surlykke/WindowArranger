function wait {
  for arg in $@; do
    while ! swaymsg "[$arg] focus" > /dev/null; do 
      sleep 0.1; 
      if [[ -n "$DEADLINE" ]] && [[ "$DEADLINE" -lt $(date +%s) ]]; then 
        echo "Window specified by $arg did not appear before timeout"
        exit 1
      fi
    done
  done  
}

function waitWithDeadline {
  export DEADLINE=$1 
  shift
  wait $@
  DEADLINE=
}




