#!/bin/bash

if [ $# -ne 2 ]
then
    echo "usage $0 prj-name dst-dir"
    exit
fi

prjName=$1
dstDir=$2

curDir=`dirname $0`
cd $curDir/../
prjHome=`pwd`

os=`uname`
isMac=0
if [ "$os" = "Darwin" ]
then
    isMac=1
fi

tmpDir=$prjHome/tmp
mkdir $tmpDir
cd $tmpDir
rm -rf $prjName

git clone git@github.com:goinbox/gobox-demo.git $prjName

cd $prjName
rm -rf .git
cd src
rm -f go.*
go mod init $prjName
cd $tmpDir

if [ $isMac -eq 1 ]
then
    sed -i '' "s/gdemo/$prjName/g" `grep -rl "gdemo" $prjName`
else
    sed -i "s/gdemo/$prjName/g" `grep -rl "gdemo" $prjName`
fi

if [ ! -d $dstDir ]
then
    mkdir -p $dstDir
fi

mv $prjName $dstDir
