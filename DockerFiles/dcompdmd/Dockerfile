FROM golang
 #:3.4
#RUN apk add --no-cache gcc libc-dev

COPY etc/dcomp/plugins/local/dmd.yaml /etc/dcomp/conf/
COPY etc/dcomp/cert/*dmd.pem /etc/dcomp/cert/
COPY bin/dcompdmd /
ENTRYPOINT ["/dcompdmd","/etc/dcomp/conf/dmd.yaml"]