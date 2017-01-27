#!/bin/bash
set -o xtrace

id=`dcomp submit -resource SlurmDocker -nnodes=2 -script "echo hello" centos:7`
echo $id > id
dcomp wait $id
dcomp cp -u $id / .

