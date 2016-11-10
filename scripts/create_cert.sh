#!/usr/bin/env bash

openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout ../etc/dcomp/keyd.pem -out ../etc/dcomp/certd.pem \
-subj '/CN=localhost:8001/O=DESY/C=DE'

openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout ../etc/dcomp/keyauth.pem -out ../etc/dcomp/certauth.pem \
-subj '/CN=localhost:800r71/O=DESY/C=DE'
