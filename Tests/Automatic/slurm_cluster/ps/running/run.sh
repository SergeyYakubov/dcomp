#!/bin/bash
set -o xtrace

id=`dcomp submit -resource SlurmDocker -nnodes=2 -script "sleep 20" none`
echo $id > id
dcomp wait -status running $id
dcomp ps -id $id > output_psid


