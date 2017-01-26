#!/usr/bin/env bash

out=`sacct -nj $1 --format=start,end,state | head -n 1`

# remove + from status (in case of staus like CANCELLED+)
echo "${out/+/}"
