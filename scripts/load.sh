#!/usr/bin/env bash

# This is a helper script for loading our plugin to expedite testing during development

PLUGIN='snap-plugin-publisher-signalfx'
TASK=${HOME}'/tasks/signalfx.yaml'

if [ ! -d ${GOPATH} ]; then echo "GOPATH may not be set correctly"; exit -1; fi

if [ -e ${GOPATH}/bin/${PLUGIN} ]; then rm ${GOPATH}/bin/${PLUGIN}; fi

if [ ! -e ${TASK} ]; then echo "Task file does not exist"; exit -1; fi

go install && \
  snaptel plugin load ~/golang/bin/${PLUGIN} && \
  snaptel task create -t ${TASK} && \
  snaptel task list
