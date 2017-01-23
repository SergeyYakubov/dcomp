#!/usr/bin/env bash

useradd $DOCKERUSER
ln -s /bin/echo /usr/bin/dockerrun
ln -s /bin/echo /usr/bin/dockercluster
ln -s /bin/echo /usr/bin/dockerexec
/dcomp/bin/dcompdmd /etc/dcomp/plugins/slurm/dmd.yaml &
/dcomp/bin/dcompclusterpd /etc/dcomp/plugins/slurm/slurm.yaml




