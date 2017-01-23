#!/bin/bash

docker run -d --name slurm -h slurmhost -p 8010:8010  -p 8009:8009 -v `pwd`/extra_start.sh:/usr/bin/extra_start.sh \
  -v $DCOMP_BASEDIR/bin:/dcomp/bin -v $DCOMP_BASEDIR/etc/dcomp/:/etc/dcomp/ -v /dcompdata \
  --add-host=daemonhost:`ip route show | grep docker0 | awk '{print \$9}'` \
  --add-host=databasehost:`ip route show | grep docker0 | awk '{print \$9}'` \
  yakser/centos_slurm

# wait for initialization
c=0
while ! `curl localhost:8009>/dev/null 2>&1`; do
    sleep 1
    ((c++))
    if [ $c -eq "10" ]; then
        echo "timeout starting slurm container"
        exit 1
    fi
done