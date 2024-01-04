#!/bin/sh

# Extract the instance number from the container name
INSTANCE_NUMBER=$(echo "$HOSTNAME" | sed 's/[^0-9]//g')

# Create a unique temporary directory for each Raft instance
mkdir -p "/tmp/$INSTANCE_NUMBER"

# Run the Raft instance
/app/raftinstance.exe "-state_dir=/tmp/$INSTANCE_NUMBER" "-raft=:$((8080+INSTANCE_NUMBER))" "-api=:$((9090+INSTANCE_NUMBER))"
