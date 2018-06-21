#!/bin/bash
SESSION_NAME="nwservices"

tmux has-session -t ${SESSION_NAME}

if [ $? != 0 ]
then
	echo "Creating ${SESSION_NAME} session..."
	# Create the session
  tmux new-session -d -s ${SESSION_NAME} -n vim
  tmux send-keys -t ${SESSION_NAME} 'cd ~/Sites/compilebox/API/' C-m './bin/compilebox' C-m

  tmux split-window -v -t ${SESSION_NAME}
  tmux send-keys -t ${SESSION_NAME}:0.1 'cd ~/Sites/testbox/' C-m './bin/testbox' C-m

  tmux a -t ${SESSION_NAME}
else
	echo "${SESSION_NAME} already running"
fi
