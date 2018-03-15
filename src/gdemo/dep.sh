#!/bin/bash

curDir=`dirname $0`
cd $curDir
curDir=`pwd`
cd $curDir/../../
prjHome=`pwd`

export GOPATH=$prjHome
export GOBIN=$prjHome/bin

cd $curDir
dep $@
