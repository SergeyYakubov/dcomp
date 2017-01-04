#!/usr/bin/env bash


pkill dcomp
dcomp daemon &
dcompestd &
dcompauthd &
dcomplocalpd &
dcompdmd /etc/dcomp/plugins/local/dmd.yaml &

