#!/usr/bin/env bash

DIR=$1
HOST=$2

openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout $DIR/keyd.pem -out $DIR/certd.pem \
-subj "/CN=$HOST/O=DESY/C=DE"

