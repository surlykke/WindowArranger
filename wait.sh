function wait {
	for arg in $@; do
		echo "Looping " swaymsg "[$arg] focus"
		while ! swaymsg "[$arg] focus" > /dev/null; do 
			sleep 0.1; 
		done
	done	
}
