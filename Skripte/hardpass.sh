#!/bin/bash

source vars
echo -e "$HARDPASS\n$HARDPASS\n" | sudo passwd $USER
