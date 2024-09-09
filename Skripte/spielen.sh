#!/bin/bash

source $(dirname $0)/vars

$(dirname $0)/simplepass.sh
while true
do
    test $(who | grep linus | wc -l) -gt 0 && break
done

$(dirname $0)/hardpass.sh
sudo at -r $(cat $ATQID)
echo "$PWD/beenden.sh" | sudo at -m now + $TIME min 2>&1 | tail -1 | cut -d' ' -f2 > $ATQID
