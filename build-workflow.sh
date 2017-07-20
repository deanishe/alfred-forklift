#!/bin/bash

here="$( cd "$( dirname "$0" )"; pwd )"


log() {
    echo "$@" > /dev/stderr
}

pushd "$here" &> /dev/null

log "cleaning ./build ..."
rm -rvf ./build

log

log "copying assets to ./build ..."

mkdir -vp ./build

cp -v icon.png ./build/
cp -v update-available.png ./build/
cp -v info.plist ./build/
cp -v README.md ./build/
cp -v LICENCE.txt ./build/

log

log "building executable(s) ..."
go build -v -o ./forklift ./forklift.go
ST_BUILD=$?
if [ "$ST_BUILD" != 0 ]; then
    log "error building executable."
    rm -rf ./build/
    popd &> /dev/null
    exit $ST_BUILD
fi

chmod 755 ./forklift
cp -v ./forklift ./build/forklift

# Get the dist filename from the executable
zipfile="$(./forklift --distname 2> /dev/null)"

log

if test -e "$zipfile"; then
    log "removing existing .alfredworkflow file ..."
    rm -rvf "$zipfile"
    log
fi

log "building .alfredworkflow file ..."
pushd ./build/ &> /dev/null
zip -v "../${zipfile}" *
ST_ZIP=$?
if [ "$ST_ZIP" != 0 ]; then
    log "error creating .alfredworkflow file."
    rm -rf ./build/
    popd &> /dev/null
    exit $ST_ZIP
fi
popd &> /dev/null

log

log "cleaning up ..."
rm -rvf ./build/

popd &> /dev/null
log "all done."

