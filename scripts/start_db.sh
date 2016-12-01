docker network create --subnet=172.18.0.0/16 mynet
docker run --net mynet --ip 172.18.0.2 --name mongodb -p 27017:27017 -d mongo
#docker run --net mynet --ip 172.18.0.3 --name dcompd -P  -d yakser/dcompd
#docker run --net mynet --ip 172.18.0.4 --name dcompestd -P  -d yakser/dcompestd
#docker run --net mynet --ip 172.18.0.5 --name dcompauthd -P  -d yakser/dcompauthd
#docker run --net mynet --ip 172.18.0.6 --name dcompdmd -P -v /dcompdata:/dcompdata -d yakser/dcompdmd
