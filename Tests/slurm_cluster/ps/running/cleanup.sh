#!/bin/bash

rm -rf output_psid output_psid_nosleep

id=`cat id`
dcomp kill $id
dcomp ps -id $id
dcomp wait $id
dcomp rm -id $id
rm id
