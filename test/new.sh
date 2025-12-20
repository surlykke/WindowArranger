#!/usr/bin/env bash
# Move everything aside
swaymsg '[title=.*] move workspace 2'

# workspace 1
swaymsg '[title=T1] move to workspace 1'
swaymsg '[title=T1] focus'
swaymsg '[title=T1] split v'
swaymsg '[title=T1] layout splith'
swaymsg '[title=T1] mark current'
swaymsg '[title=T4] move to mark current'
swaymsg '[title=T3] move to mark current'
swaymsg '[title=T2] move to mark current'
swaymsg 'unmark current'
swaymsg '[title=T1] focus'
swaymsg '[title=T1] split v'
swaymsg '[title=T1] layout splitv'
swaymsg '[title=T1] mark current'
swaymsg '[title=T2] move to mark current'
swaymsg 'unmark current'
swaymsg '[title=T3] focus'
swaymsg '[title=T3] split v'
swaymsg '[title=T3] layout splith'
swaymsg '[title=T3] mark current'
swaymsg '[title=T4] move to mark current'
swaymsg 'unmark current'

swaymsg move workspace to eDP-1

