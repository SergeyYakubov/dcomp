#!/bin/bash
set -o xtrace
set -e


id=`dcomp submit -ncpus=2 -script "sleep 100" none`
dcomp wait -status running $id
dcomp kill $id
dcomp wait -status cancelled $id
dcomp rm -id $id
docker exec -i slurm test ! -d /dcompdata/$id
