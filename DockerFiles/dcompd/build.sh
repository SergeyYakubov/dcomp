#!/usr/bin/env bash
cd ../..
docker build -t yakser/dcompd -f DockerFiles/dcompd/Dockerfile .
cd -
