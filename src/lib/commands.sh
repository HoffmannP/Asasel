#!/bin/bash

function findjob {
    key=$1
    for jobid in $(atq | tr -s ' ' "\t" | cut -f1)
    do
        if [[ $(at -c $jobid | tail -2 | head -1 | rg -c "#$key") -eq 1 ]]
        then
            echo $jobid
            return
        fi
    done
}

case "$1" in
    chpwd)
        user="$2"
        pass="$3"
        # echo 'echo -en "'"$pass"'\n'"$pass"'\n" | passwd '"$user"
        echo -en "$pass\n$pass\n" | passwd $user
        ;;
    chmod)
        mode="$2"
        shift; shift; cmds="$@"
        for cmd in $cmds
        do
            # echo "chmod $mode $cmd"
            chmod $mode $cmd
        done
        ;;
    setTimelimit)
        user="$2"
        time="$3"
        # echo 'echo "killall -u '"$user"' #timelimit" | at now + '"$time"' minutes'
        echo "killall -u $user #timelimit" | at now + $time minutes
        ;;
    showTimelimit)
        job=$(findjob 'timelimit')
        if [[ -n "$job" ]]
        then
            atq $job | tr -s ' ' "\t" | cut -f5
        else
            echo "No job found" >&2
        fi
        ;;
    rmTimelimit)
        job=$(findjob 'timelimit')
        if [[ -n "$job" ]]
        then
            atrm $job
        else
            echo "No job found" >&2
        fi
        ;;
esac