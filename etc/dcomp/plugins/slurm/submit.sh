#!/usr/bin/env bash

set -e

out=`sbatch $1`

#extract job id
echo $out | cut -d ' ' -f4