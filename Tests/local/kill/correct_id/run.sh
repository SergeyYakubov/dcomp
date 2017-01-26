#!/bin/bash
set -o xtrace

id=`dcomp submit -local -nnodes 1 -script "sleep 100" centos:7`
echo $id > id
dcomp wait -status running $id
dcomp kill $id
ec=$?
dcomp ps -id $id
exit $ec
