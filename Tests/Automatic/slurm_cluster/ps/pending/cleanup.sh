#!/bin/bash

rm -rf output_psid output_psid_nosleep

for i in `seq 1 3`;
do
id=`cat id$i`
dcomp kill $id
dcomp wait $id
dcomp rm -id $id
rm id$i
done
