cmake_minimum_required (VERSION 2.8)

project (DCOMP)

option (DISABLE_CGO "Disable CGo" OFF)

if (DISABLE_CGO)
    set (GO_ENABLED "CGO_ENABLED=0")
endif()




set (DCOMP_CMAKE_DIR "${DCOMP_SOURCE_DIR}/CMake")

#set (CMAKE_MODULE_PATH ${CMAKE_MODULE_PATH} ${DCOMP_CMAKE_DIR})

add_custom_target(version ALL
  COMMAND ${CMAKE_COMMAND} -P ${DCOMP_CMAKE_DIR}/PrepareVersion.cmake
  VERBATIM)

add_custom_target(dcomp ALL
				  COMMAND  ${GO_ENABLED} go build ${GO_OPTS} -o bin/dcomp -tags release ./dcomp/main
				  WORKING_DIRECTORY ${DCOMP_SOURCE_DIR}
				  VERBATIM)

add_custom_target(dcompestd ALL
        				  COMMAND ${GO_ENABLED} go build ${GO_OPTS} -o bin/dcompestd ./dcomp/estimator/main
        				  WORKING_DIRECTORY ${DCOMP_SOURCE_DIR}
        				  VERBATIM)

add_custom_target(dcompauthd ALL
        				  COMMAND ${GO_ENABLED} go build ${GO_OPTS} -o bin/dcompauthd ./dcomp/authorizer/main
        				  WORKING_DIRECTORY ${DCOMP_SOURCE_DIR}
        				  VERBATIM)


add_custom_target(dcomplocalpd ALL
        				  COMMAND ${GO_ENABLED} go build ${GO_OPTS} -o bin/dcomplocalpd ./dcomp/resources/local/main
        				  WORKING_DIRECTORY ${DCOMP_SOURCE_DIR}
        				  VERBATIM)

add_custom_target(dcompclusterpd ALL
        				  COMMAND ${GO_ENABLED} go build ${GO_OPTS} -o bin/dcompclusterpd ./dcomp/resources/cluster/main
        				  WORKING_DIRECTORY ${DCOMP_SOURCE_DIR}
        				  VERBATIM)



add_custom_target(dcompdmd ALL
        				  go build ${GO_OPTS}  -o bin/dcompdmd ./dcomp/datamanager/main
        				  WORKING_DIRECTORY ${DCOMP_SOURCE_DIR}
        				  VERBATIM)



add_dependencies (dcomp version)
add_dependencies (dcompestd version)
add_dependencies (dcompauthd version)
add_dependencies (dcomplocalpd version)
add_dependencies (dcompdmd version)
