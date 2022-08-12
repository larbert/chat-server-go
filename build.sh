#!/bin/bash

# 自动编译项目

# 获取包名
package_name=`go list`

# 生成位置
target=./target

rm -rf $target/*

mkdir $target
echo 创建目录$target

mkdir $target/server
echo 创建目录$target/server
mkdir $target/server/mac
echo 创建目录$target/mac
mkdir $target/server/linux
echo 创建目录$target/linux
mkdir $target/server/windows
echo 创建目录$target/windows

mkdir $target/client
echo 创建目录$target/client
mkdir $target/client/mac
echo 创建目录$target/mac
mkdir $target/client/linux
echo 创建目录$target/linux
mkdir $target/client/windows
echo 创建目录$target/windows

# 编译server
# mac
go build
mv $package_name $target/server/mac
echo 生成$target/server/mac/$package_name
# liunx
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build
mv $package_name $target/server/linux
echo 生成$target/server/linux/$package_name
# windows
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build
mv $package_name.exe $target/server/windows
echo 生成$target/server/windows/$package_name.exe

# 编译client

cd client

client_package=`go list`
client_name=${client_package#*/}

# mac
go build
mv $client_name ../$target/client/mac
echo 生成$target/client/mac/$client_name
# liunx
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build
mv $client_name ../$target/client/linux
echo 生成$target/client/linux/$client_name
# windows
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build
mv $client_name.exe ../$target/client/windows
echo 生成$target/client/windows/$client_name.exe