#!/bin/bash
set -o xtrace

id1=`dcomp submit -nnodes=2 -script "sleep 20" none`
id2=`dcomp submit -nnodes=2 -script "sleep 20" none`
id3=`dcomp submit -nnodes=2 -script "sleep 20" none`
dcomp wait -wait-changes  $id3
dcomp ps -id $id3 > output_psid


