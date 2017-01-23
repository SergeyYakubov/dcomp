#!/bin/bash
set -o xtrace

id=`dcomp submit -nnodes=2 -script "srun -n 2 echo hello" none`
sleep 3
dcomp cp -u $id / .
rename 's/-(.+)/.out/' *.out

