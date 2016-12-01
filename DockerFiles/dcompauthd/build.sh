#!/usr/bin/env bash
cd ../..
docker build -t yakser/dcompauthd -f DockerFiles/dcompauthd/Dockerfile .
cd -
