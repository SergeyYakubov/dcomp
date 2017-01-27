#!/bin/bash
set -o xtrace

id=`dcomp submit -resource SlurmDocker -nnodes=2 -script "srun -n 2 echo hello" none`
echo $id > id
dcomp wait $id
dcomp ps -id $id > output_psid
dcomp ps -id $id -log > output_log

