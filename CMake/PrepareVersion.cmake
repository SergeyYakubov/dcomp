# set path for modules


execute_process(COMMAND git describe --tags --dirty OUTPUT_VARIABLE VERSION)
string(STRIP ${VERSION} VERSION)

string(REGEX REPLACE "^v([0-9]+)\\..*" "\\1" VERSION_MAJOR "${VERSION}")
string(REGEX REPLACE "^v[0-9]+\\.([0-9]+).*" "\\1" VERSION_MINOR "${VERSION}")
string(REGEX REPLACE "^v[0-9]+\\.[0-9]+\\.([0-9]+).*" "\\1" VERSION_PATCH "${VERSION}")
string(REGEX REPLACE "^v[0-9]+\\.[0-9]+\\.[0-9]+(.*)" "\\1" VERSION_SHA1 "${VERSION}")
set(VERSION_SHORT "${VERSION_MAJOR}.${VERSION_MINOR}.${VERSION_PATCH}")

string(TIMESTAMP TIMESTAMP "%H:%M:%S %d.%m.%Y UTC" UTC)

configure_file(${CMAKE_SOURCE_DIR}/dcomp/version/version_lib.go.in ${CMAKE_SOURCE_DIR}/dcomp/version/version_lib.go)

