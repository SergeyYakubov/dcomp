#!/bin/bash
#set -o xtrace

id=$(dcomp submit -resource local -nnodes=1 -upload `pwd`/script.sh:/data/script.sh -script /data/script.sh centos:7)
echo $id > id
dcomp wait $id
dcomp ps -log -id $id
