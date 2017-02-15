#!/bin/bash

id=`cat id`
dcomp rm -id $id
rm id

id=`cat id2`
dcomp rm -id $id
rm id2
