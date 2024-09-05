#!/bin/bash

source vars
LOGFILEDIR="/home/$USER/.minecraft"

for file in $(ls -rt "$LOGFILEDIR/launcher_log*.txt")
do
    echo -n "On "
    head -1 $file | cut -d':' -f2-3
    echo -n "Off"
    tail -1 $file | cut -d':' -f2-3
done
ps -eo lstart,cmd | rg minecraft | rg java | while read instance; do
    echo $instance  | tr -s ' ' "\t" | cut -f -4 | tr -s "\t" " " | date +"On  %Y-%m-%d %H:%M" -f -
done