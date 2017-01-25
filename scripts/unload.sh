#!/usr/bin/env bash

# This is a helper script for unloading our plugin to expedite testing during development

snaptel task list | tail -1 | awk '{ print $1 }' | xargs snaptel task stop

snaptel task list | tail -1 | awk '{ print $1 }' | xargs snaptel task remove

snaptel plugin unload publisher signalfx 1

if [ -e /tmp/signalfx-debug.log ];then sudo rm -f /tmp/signalfx-debug.log; fi
