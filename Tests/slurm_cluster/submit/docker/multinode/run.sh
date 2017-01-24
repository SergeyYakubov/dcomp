#!/bin/bash
set -o xtrace

id=`dcomp submit -nnodes=2 -script "echo hello" centos:7`
sleep 3
dcomp cp -u $id / .

