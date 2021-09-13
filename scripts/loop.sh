#!/bin/sh

dir=$(dirname $0)

while [[ 1 ]]; do
    echo
    echo
    date
    git pull
    $dir/screenshot.sh --pause 10s "$@"
    $dir/commit.sh
    sleep 1800
done
