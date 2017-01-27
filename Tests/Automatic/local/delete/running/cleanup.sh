#!/bin/bash

id=`cat id`
dcomp kill $id
dcomp wait $id
dcomp rm -id $id
rm id
