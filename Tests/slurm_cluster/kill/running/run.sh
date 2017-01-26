#!/bin/bash
set -o xtrace
set -e


id=`dcomp submit -ncpus=2 -script "sleep 100" none`
echo $id > id
dcomp wait -status running $id
dcomp ps -id $id > output_running
dcomp kill $id
dcomp wait -status cancelled $id
dcomp ps -id $id > output_cancelled
