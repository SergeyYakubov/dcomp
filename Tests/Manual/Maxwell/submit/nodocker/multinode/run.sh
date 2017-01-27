#!/bin/bash
set -o xtrace

id=`dcomp submit -resource maxwell -nnodes=2 -script "srun -n 2 echo hello" none`
echo $id > id
dcomp wait $id
dcomp cp -u $id / .

