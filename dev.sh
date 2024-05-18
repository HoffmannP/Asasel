#!/bin/bash

set -e

id=$(kitty @ new-window --cwd $(pwd))
kitty @ send-text -m id:$id '(cd backend; ./dev.sh)\n'
kitty @ goto-layout vertical

(cd frontend; ./dev.sh)