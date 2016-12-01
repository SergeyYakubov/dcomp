FROM scratch

COPY etc/dcomp/conf/dcompauthd.yaml /etc/dcomp/conf/
COPY etc/dcomp/cert/certauth.pem /etc/dcomp/cert/
COPY etc/dcomp/cert/keyauth.pem /etc/dcomp/cert/
COPY bin/dcompauthd /
ENTRYPOINT ["/dcompauthd"]