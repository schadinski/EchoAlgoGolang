#!/bin/bash

# Find node processes and kill them
# PIDS=`ps -e | grep nodeMain | cut -d" " -f2`
# for pid in $PIDS
# do 
#     `kill -9 $pid`
#     echo "Killed prozess $pid"
# done
pkill -f nodeMain
echo "All nodes killed"

# Check for args to kill logger
if [ -n $1 ]; then
    EXPECTED_ARG="l"
    GIVEN_ARG=$1
    if [[ $GIVEN_ARG == $EXPECTED_ARG ]]; then
        # LOGGER=`ps -e | grep loggerMain | cut -d" " -f2`
        # `kill -9 $LOGGER`
        pkill -f nodeLogger
        echo "logger killed"
    fi
fi