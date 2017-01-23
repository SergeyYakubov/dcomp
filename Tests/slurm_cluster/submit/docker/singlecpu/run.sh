#!/bin/bash
set -o xtrace

id=`dcomp submit -ncpus=1 -script "echo hello" centos:7`
sleep 3
dcomp cp -u $id / .
rename 's/-(.+)/.out/' *.out
cat slurm.out

