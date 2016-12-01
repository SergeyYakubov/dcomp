#!/usr/bin/env bash
cd ../..
docker build -t yakser/dcompestd -f DockerFiles/dcompestd/Dockerfile .
cd -
