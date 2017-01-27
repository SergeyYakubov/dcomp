#!/bin/bash

id=`cat id`
dcomp rm -id $id
rm id

rm output_cancelled
rm output_running