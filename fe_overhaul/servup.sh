#!/bin/bash
SESSION_NAME="nwstorybook"

tmux has-session -t ${SESSION_NAME}

if [ $? != 0 ]
then
	echo "Creating ${SESSION_NAME} session..."
	# Create the session
  tmux new-session -d -s ${SESSION_NAME} -n vim
  tmux send-keys -t ${SESSION_NAME} 'npm run storybook' C-m
else
	echo "${SESSION_NAME} already running"
fi
