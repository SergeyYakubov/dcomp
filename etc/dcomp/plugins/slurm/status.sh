#!/usr/bin/env bash

sacct -nj $1 --format=start,end,state | head -n 1
