#!/bin/bash
set -o xtrace

id=$(dcomp submit -resource local -nnodes=1 -upload `pwd`/script.sh:/data/aaa/script.sh -script /data/aaa/script.sh centos:7)
echo $id > id
dcomp wait $id

id2=$(dcomp submit -resource local -nnodes=1 -mount $id/data:/data -script /data/aaa/script.sh centos:7)
echo $id2 > id2
dcomp wait $id2

dcomp ps -log -id $id2
