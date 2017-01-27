#!/bin/bash
set -o xtrace

id=`dcomp submit -resource local -nnodes 1 -script "echo hello" centos:7`
echo $id > id
dcomp wait $id
dcomp rm -id $id
test ! -d /dcompdata/$id