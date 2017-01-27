#!/bin/bash
set -o xtrace

id=`dcomp submit -resource local -nnodes 1 -script "sleep 10" centos:7`
echo $id > id
dcomp wait -status running $id
dcomp rm -id $id
