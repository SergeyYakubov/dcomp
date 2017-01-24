#!/bin/bash
#set -o xtrace

id=$(dcomp submit -local -nnodes=1 -upload `pwd`/script.sh:/data -script /data/script.sh centos:7)
dcomp wait $id
dcomp ps -log -id $id
