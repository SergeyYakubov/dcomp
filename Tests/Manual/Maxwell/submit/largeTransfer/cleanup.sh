#!/bin/bash

rm -rf job.sh job.log

id=`cat id`
dcomp rm -id $id
rm id