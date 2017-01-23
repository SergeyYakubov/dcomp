#!/bin/bash
set -o xtrace

id=`dcomp submit -local -nnodes 1 -script "echo hello" centos:7`
dcomp wait $id
dcomp ps -log -id $id
