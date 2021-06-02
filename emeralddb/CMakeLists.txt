cmake_minimum_required(VERSION 3.17)
project(emeralddb)

set(CMAKE_CXX_STANDARD 98)

include_directories(src/include)

file(GLOB_RECURSE SOURCE_FILES "src/*.cpp" "src/*.hpp")
add_executable(emeralddb ${SOURCE_FILES})

# 设置boost使用静态链接
set(Boost_USE_STATIC_LIBS ON)
# 设置需要的boost 组件
set(BOOST_COMPONENTS program_options system thread)

find_package(Boost REQUIRED COMPONENTS ${BOOST_COMPONENTS})
find_package(bson-1.0 REQUIRED)

target_link_libraries(emeralddb ${Boost_LIBRARIES})
target_link_libraries(emeralddb mongo::bson_static)

