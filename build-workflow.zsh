#!/usr/bin/env zsh

set -e

here="$( cd "$( dirname "$0" )"; pwd )"
source "${here}/alfred_env.sh"

build=true
devmode=true
verbose=

usage() {
    cat <<EOS
build-workflow.sh [options]

Usage:
    build-workflow.sh [-x] [-d] [-v]
    build-workflow.sh -h

Options:
    -d      Distribution. Also build .alfredworkflow file.
    -x      Don't build executable.
    -v      Verbose.
    -h      Show this help message and exit.
EOS
}

while getopts ":dhvx" opt; do
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
    x)
      build=false
      ;;
    \?)
      log "Invalid option: -$OPTARG"
      exit 1
      ;;
  esac
done
shift $((OPTIND-1))


log() {
    echo "$@" > /dev/stderr
}

pushd "$here" &> /dev/null

test -d "build" && {
    log "cleaning ./build ..."
    rm $verbose -rf ./build
    log
}


log "copying assets to ./build ..."

mkdir $verbose -p ./build

ln $verbose icon.png ./build/
ln $verbose update-available.png ./build/
ln $verbose info.plist ./build/
ln $verbose README.md ./build/
ln $verbose LICENCE.txt ./build/
log

$build && {
    log "building executable(s) ..."
    go build -v -o ./forklift ./forklift.go
    ST_BUILD=$?
    if [ "$ST_BUILD" != 0 ]; then
        log "error building executable."
        rm $verbose -rf ./build/
        popd &> /dev/null
        exit $ST_BUILD
    fi
}

chmod 755 ./forklift
cp $verbose ./forklift ./build/forklift

# Get the dist filename from the executable
zipfile="$(./forklift --distname 2> /dev/null)"

if test -e "$zipfile"; then
    log "removing existing .alfredworkflow file ..."
    rm $verbose -rf "$zipfile"
    log
fi

$devmode || {
    log "building $zipfile ..."
    pushd ./build/ &> /dev/null
    zip $verbose "../${zipfile}" *
    ST_ZIP=$?
    if [ "$ST_ZIP" != 0 ]; then
        log "error creating .alfredworkflow file."
        rm $verbose -rf ./build/
        popd &> /dev/null
        exit $ST_ZIP
    fi
    popd &> /dev/null
    log
}


popd &> /dev/null
log "all done."

