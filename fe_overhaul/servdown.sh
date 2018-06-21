#!/bin/bash
SESSION_NAME="nwstorybook"


tmux has-session -t ${SESSION_NAME}

if [ $? != 0 ]
then
	echo "${SESSION_NAME} not running"
else
	echo "taking down ${SESSION_NAME}"
	tmux kill-session -t ${SESSION_NAME}
fi
