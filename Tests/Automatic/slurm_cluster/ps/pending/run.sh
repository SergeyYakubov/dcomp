#!/bin/bash
set -o xtrace

id1=`dcomp submit -resource SlurmDocker -nnodes=2 -script "sleep 20" none`
id2=`dcomp submit -resource SlurmDocker -nnodes=2 -script "sleep 20" none`
id3=`dcomp submit -resource SlurmDocker -nnodes=2 -script "sleep 20" none`
echo $id1 > id1
echo $id2 > id2
echo $id3 > id3
dcomp wait -wait-changes  $id3
dcomp ps -id $id3 > output_psid


