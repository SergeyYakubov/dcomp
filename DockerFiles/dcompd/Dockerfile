FROM scratch
#FROM alpine:3.4
#RUN apk add --no-cache ca-certificates apache2-utils

COPY etc/dcomp/conf/dcompd.yaml /etc/dcomp/conf/
COPY etc/dcomp/cert/certd.pem /etc/dcomp/cert/
COPY etc/dcomp/cert/keyd.pem /etc/dcomp/cert/

COPY bin/dcomp /dcomp
ENTRYPOINT ["/dcomp","daemon"]