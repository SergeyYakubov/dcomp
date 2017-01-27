#!/bin/bash

rm -rf output_log output_psid

id=`cat id`
dcomp rm -id $id
rm id
