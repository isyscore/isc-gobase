cmake_minimum_required(VERSION 3.21)
project(gojvm)

set(CMAKE_CXX_STANDARD 23)

include_directories(/usr/local/include)
link_libraries(/usr/local/lib/libjvm.dylib)

set(SRCS
    gojvm_c.c
)

set(HDRS
    gojvm_c.h
    gojvm_wrap.h
)

add_library(gojvm ${SRCS} ${HDRS})