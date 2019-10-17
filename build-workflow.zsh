#!/usr/bin/env zsh

set -e

here="${${(%):-%x}:A:h}"
source "${here}/.env"

devmode=true
verbose=

log() {
    echo "$@" > /dev/stderr
}

usage() {
    cat <<EOS
build-workflow.zsh [options]

Usage:
    build-workflow.zsh [-d] [-v]
    build-workflow.zsh -h

Options:
    -d      also build .alfredworkflow file
    -v      verbose
    -h      show this help message and exit
EOS
}

while getopts ":dhv" opt; do
  case $opt in
    d)
      devmode=false
      ;;
    h)
      usage
      exit 0
      ;;
    v)
      verbose=-v
      ;;
    \?)
      log "invalid option: -$OPTARG"
      exit 1
      ;;
  esac
done
shift $((OPTIND-1))


cd "$here"

test -d "build" && {
    log "cleaning ./build ..."
    command rm $verbose -rf ./build/*
    log
} || true


log "copying assets to ./build ..."

mkdir $verbose -p ./build

cd ./build
ln -s $verbose ../*.png .
ln -s $verbose ../info.plist .
ln -s $verbose ../README.md .
ln -s $verbose ../LICENCE.txt .
cd -
log


log "building executable(s) ..."
go build -v -o ./build/forklift .
log


$devmode || {
    # Get the dist filename from the executable
    zipfile="ForkLift-Favourites-${alfred_workflow_version}.alfredworkflow"

    test -f "./dist/$zipfile" && {
        log "removing existing .alfredworkflow file ..."
        rm $verbose -f "./dist/$zipfile"
        log
    } || true

    log "building $zipfile ..."
    mkdir $verbose -p ./dist
    cd ./build
    zip $verbose "../dist/${zipfile}" *
    log
}

log "all done"

