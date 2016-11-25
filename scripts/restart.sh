#!/usr/bin/env bash


pkill dcomp
dcomp daemon &
dcompestd &
dcompauthd &
dcomplocalpd &
dcompdmd /etc/dcomp/plugins/local/local_dmd.yaml &

