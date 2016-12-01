#!/usr/bin/env bash
cd ../..
docker build -t yakser/dcompdmd -f DockerFiles/dcompdmd/Dockerfile .
cd -
