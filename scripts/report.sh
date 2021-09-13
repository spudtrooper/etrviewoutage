#!/bin/sh
#
# Creates report.
#
set -e

dir=$(dirname $0)

go run $dir/../report.go "$@"