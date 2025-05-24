#!/bin/bash

source $(dirname $0)/vars

LOGFILEDIR="/home/$USER/.minecraft/logs"

function printdur {
    echo "$(printf %3d $[$1 / 60]):$(printf %02d $[$1 % 60])"
}

total=0
first=$[$[$[$(date +%u) + 5 ] % 7] + 1]
for ((i=$first; i>=0; i--))
do
    then=$[$(date +%s) - i * 60 * 60 * 24]
    day=$(date -d @$then +%F)
    for file in $(ls $LOGFILEDIR/$day-*.log.gz 2>/dev/null)
    do
        start=$(zcat $file | head -1 | sed 's/^.//;s/:..].*//')
        end=$(zcat $file | tail -1 | sed 's/^.//;s/:..].*//')
        duration=$[($(date -d $end +%s) - $(date -d $start +%s)) / 60]
        echo -e "$(date -d @$then +"%a %F") $start - $end $(printdur $duration)"
        total=$[$total + $duration]
    done
done

file=$LOGFILEDIR/latest.log
then=$(stat -c%Y $file)
day=$(date -d @$then +%F)
start=$(cat $file | head -1 | sed 's/^.//;s/:..].*//')
end=$(cat $file | tail -1 | sed 's/^.//;s/:..].*//')
duration=$[($(date -d $end +%s) - $(date -d $start +%s)) / 60]
echo -e "$(date -d @$then +"%a %F") $start - $end $(printdur $duration)"
total=$[$total + $duration]
echo "----------------------------------"
echo "Total in week $(date -d @$then +%V)            $(printdur $total)"


# Die letzten 10 mal
# for file in $(ls -rt $LOGFILEDIR/../launcher_log*.txt)
# do
#     echo -n "On "
#     head -1 $file | cut -d':' -f2-3
#     echo -n "Off"
#     tail -1 $file | cut -d':' -f2-3
# done

ps -eo lstart,cmd | rg minecraft | rg java | while read instance; do
    start=$(echo $instance  | tr -s ' ' "\t" | cut -f -4 | tr -s "\t" " " | date +%s -f -)
    duration=$[$(date +%s) - $start]
    echo "Running since $(date -d @$start +%H:%M)         $(printdur $duration)"

done