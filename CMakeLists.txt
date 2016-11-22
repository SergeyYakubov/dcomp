cmake_minimum_required (VERSION 2.8)

project (DCOMP)


set (DCOMP_CMAKE_DIR "${DCOMP_SOURCE_DIR}/CMake")

#set (CMAKE_MODULE_PATH ${CMAKE_MODULE_PATH} ${DCOMP_CMAKE_DIR})

add_custom_target(make_version ALL
  COMMAND ${CMAKE_COMMAND} -P ${DCOMP_CMAKE_DIR}/PrepareVersion.cmake
  VERBATIM)

add_custom_target(dcomp ALL
				  COMMAND go build -o bin/dcomp -tags release ./dcomp/main
				  WORKING_DIRECTORY ${DCOMP_SOURCE_DIR}
				  VERBATIM)

add_custom_target(dcompestd ALL
        				  COMMAND go build -o bin/dcompestd ./dcomp/estimator/main
        				  WORKING_DIRECTORY ${DCOMP_SOURCE_DIR}
        				  VERBATIM)

add_custom_target(dcompauthd ALL
        				  COMMAND go build -o bin/dcompauthd ./dcomp/authorizer/main
        				  WORKING_DIRECTORY ${DCOMP_SOURCE_DIR}
        				  VERBATIM)


add_custom_target(dcomplocalpd ALL
        				  COMMAND go build -o bin/dcomplocalpd ./dcomp/resources/local/main
        				  WORKING_DIRECTORY ${DCOMP_SOURCE_DIR}
        				  VERBATIM)

add_custom_target(dcompdmd ALL
        				  COMMAND go build -o bin/dcompdmd ./dcomp/datamanager/main
        				  WORKING_DIRECTORY ${DCOMP_SOURCE_DIR}
        				  VERBATIM)



add_dependencies (dcomp make_version)
add_dependencies (dcompestd make_version)
add_dependencies (dcompauthd make_version)
add_dependencies (dcomplocalpd make_version)
add_dependencies (dcompdmd make_version)
