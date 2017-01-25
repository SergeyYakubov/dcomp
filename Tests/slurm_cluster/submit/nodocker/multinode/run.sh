#!/bin/bash
set -o xtrace

id=`dcomp submit -nnodes=2 -script "srun -n 2 echo hello" none`
dcomp wait $id
dcomp cp -u $id / .

