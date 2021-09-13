#!/bin/sh
#
# Tests animate.html locally
#

dir=$(dirname $0)

$dir/report.sh

echo
echo
echo Visit: http://localhost:8000/html/animate.html
echo
echo
pushd data
python -m SimpleHTTPServer
popd