#export GOPATH:=$(CURDIR)/Godeps/_workspace:$(GOPATH)

#LIBDIR=${DESTDIR}/lib/systemd/system
#BINDIR=${DESTDIR}/usr/lib/docker/

#PACKAGES = $(shell find ./ -type d -not -path '*/\.*')
SOURCEDIR=./source
PACKAGES=cmd

.PHONY: test-cover-html

all:
	go build  -o dcomp ${SOURCEDIR}/cmd

test:
	go test ${SOURCEDIR}/${PACKAGES}

coverage:
	echo "mode: count" > coverage-all.out
	$(foreach pkg,${SOURCEDIR}/$(PACKAGES),\
		go test -coverprofile=coverage.out  $(pkg);\
		tail -n +2 coverage.out >> coverage-all.out;)
	go tool cover -html=coverage-all.out -o coverage.html
	rm -rf coverage-all.out coverage.out

clean:
	rm dcomp *~
