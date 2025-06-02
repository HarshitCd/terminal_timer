#!/bin/bash

SESSION_NAME="timers"

tmux has-session -t $SESSION_NAME 2>/dev/null

if [ $? != 0 ]; then
    tmux new-session -d -s $SESSION_NAME "terminal_timer $1"
else
    tmux new-window -t $SESSION_NAME "terminal_timer $1"
fi

NEW_WINDOW_INDEX=$(tmux list-windows -t "$SESSION_NAME" -F "#{window_index}" | tail -n 1)
tmux switch-client -t $SESSION_NAME
tmux select-window -t "$SESSION_NAME:$NEW_WINDOW_INDEX"


