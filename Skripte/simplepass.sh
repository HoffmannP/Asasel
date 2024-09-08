#!/bin/bash

source vars
echo -e "$SIMPLEPASS\n$SIMPLEPASS\n" | sudo passwd -q $USER
