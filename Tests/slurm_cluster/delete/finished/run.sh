#!/bin/bash
set -o xtrace
set -e


id=`dcomp submit -ncpus=2 -script "echo hello" none`
dcomp wait $id
dcomp rm -id $id
docker exec -i slurm test ! -d /dcompdata/$id
