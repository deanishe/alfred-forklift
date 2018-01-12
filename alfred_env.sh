#!/bin/bash

here="$( cd "$( dirname "$0" )"; pwd )"
ip="${here}/info.plist"
pl=/usr/libexec/PlistBuddy

export alfred_workflow_bundleid="$( $pl -c "Print :bundleid" "$ip" )"
export alfred_workflow_name="$( $pl -c "Print :name" "$ip"  )"
export alfred_workflow_version="$( $pl -c "Print :version" "$ip"  )"
