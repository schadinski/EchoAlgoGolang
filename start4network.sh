#!/bin/bash

# Logger Port 1042
LOGGER="127.0.0.1:1042"

# Start logger process
cd Logger
go build loggerMain.go logger.go msg.go
gnome-terminal -e "./loggerMain" &
echo "logger started" 
# Start node processes
cd ../Nodes 
go build nodeMain.go node.go msg.go

gnome-terminal -e "./nodeMain "127.0.0.1:1053" $LOGGER "1" "127.0.0.1:1063" "127.0.0.1:1073" "127.0.0.1:1083" " 
gnome-terminal -e "./nodeMain "127.0.0.1:1063" $LOGGER "2" "127.0.0.1:1053" "127.0.0.1:1073" "127.0.0.1:1083" " 
gnome-terminal -e "./nodeMain "127.0.0.1:1073" $LOGGER "3" "127.0.0.1:1053" "127.0.0.1:1063" "127.0.0.1:1083" " 
gnome-terminal -e "./nodeMain "127.0.0.1:1083" $LOGGER "4" "127.0.0.1:1053" "127.0.0.1:1063" "127.0.0.1:1073" " 