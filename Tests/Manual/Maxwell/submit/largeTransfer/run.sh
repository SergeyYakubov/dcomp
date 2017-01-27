#!/bin/bash
set -o xtrace
set -e 
if [ ! -f bigfile ]; then
	dd if=/dev/zero of=bigfile bs=1G count=1
fi

id=`time dcomp submit -upload bigfile:/data -resource maxwell -nnodes=1 -script "echo hello" none`
echo $id > id
dcomp wait $id
dcomp ls $id /data

