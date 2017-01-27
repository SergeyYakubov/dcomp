#!/bin/bash
set -o xtrace

id=`dcomp submit -resource maxwell -nnodes=3 -script "mpirun hostname" yakser/centos_mpi`
echo $id > id
dcomp wait $id
dcomp cp -u $id / .

