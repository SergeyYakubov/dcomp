#!/bin/bash
set -o xtrace

id=`dcomp submit -nnodes=2 -script "sleep 20" none`
sleep 3
dcomp ps -id $id > output_psid


