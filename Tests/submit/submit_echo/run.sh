#!/bin/bash
set -o xtrace

id=`dcomp submit -local -script "echo hello" centos:7`
dcomp wait $id
dcomp ps -log -id $id
