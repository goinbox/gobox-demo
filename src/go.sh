#!/bin/bash

curDir=`dirname $0`
cd $curDir/../
prjHome=`pwd`

export GOPATH=$prjHome
export GOBIN=$prjHome/bin

go $@
