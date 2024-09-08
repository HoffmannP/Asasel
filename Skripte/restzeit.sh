#!/bin/bash

source $(dirname $0)/vars

END=$(sudo at -l | rg $(cat $ATQID) | tr -s ' ' "\t" | cut -f2-5 | tr "\t" ' ' | date +"%Y-%m-%d %H:%M" -f -)
echo $END
python -c 'from datetime import datetime; print( datetime.strptime("'"$END"'", "%Y-%m-%d %H:%M") - datetime.now() )'
