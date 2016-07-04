#!/bin/bash

source /home/yakubov/.bashrc

mapfile -t PACKAGES < <( find ./$1 -type d -not -path '*/\.*' )
#PACKAGES=./utils


echo "mode: count" > coverage-all.out
for pkg in ${PACKAGES[@]}
do
	echo $pkg
	go test -coverprofile=coverage.out  $pkg
	tail -n +2 coverage.out >> coverage-all.out
done
go tool cover -html=coverage-all.out -o coverage.html
rm -rf coverage-all.out coverage.out
firefox ./coverage.html &

