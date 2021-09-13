#!/bin/sh
#
# Creates report and publishes files to github.io directory.
#
set -e

dir=$(dirname $0)
githubio=../spudtrooper.github.io

$dir/report.sh

pushd $githubio && \
  git reset --hard HEAD &&\
  git pull && \
  popd


mkdir -p $githubio/ida/html
cp -R data/html/* $githubio/ida/html

mkdir -p $githubio/ida/screenshots
cp -R data/screenshots/* $githubio/ida/screenshots

pushd $githubio && \
  ./scripts/commit.sh && \
  popd
