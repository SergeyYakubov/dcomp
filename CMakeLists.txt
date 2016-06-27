cmake_minimum_required (VERSION 2.8)

project (DCOMP)


set (DCOMP_CMAKE_DIR "${DCOMP_SOURCE_DIR}/CMake")

#set (CMAKE_MODULE_PATH ${CMAKE_MODULE_PATH} ${DCOMP_CMAKE_DIR})

add_custom_target(make_version ALL
  COMMAND ${CMAKE_COMMAND} -P ${DCOMP_CMAKE_DIR}/PrepareVersion.cmake
  VERBATIM)

add_custom_target(dcomp ALL
				  COMMAND go build -o dcomp ./source/main
				  WORKING_DIRECTORY ${DCOMP_SOURCE_DIR}
				  VERBATIM)
					
add_dependencies (dcomp make_version)
