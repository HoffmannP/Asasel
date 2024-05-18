#!/bin/bash

set -e


export SERVER_NAME=$(git remote get-url origin | cut -d@ -f2 | sed -r 's/:/:5050\//;s/.git//')
NAME=$(echo $SERVER_NAME | rev | cut -d\/ -f1 | rev)
NAMESPACE="default"
KUBECTL="kubectl"

if [[ "$1" != "--nobuild" ]]
then
    (
        cd frontend
        npm run build
    )
    buildah build -t "${SERVER_NAME}:latest" .
else
    shift
fi

VERSION="$(date '+%Y-%V-%a')"
buildah tag "$SERVER_NAME:latest" "$SERVER_NAME:$VERSION"
buildah push "$SERVER_NAME:latest"
buildah push "$SERVER_NAME:$VERSION"

if [[ "$1" != "--nodeploy" ]]
then
    $KUBECTL -n $NAMESPACE rollout restart deploy $NAME
else
    shift
fi

