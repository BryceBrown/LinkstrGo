#!/bin/bash
echo "Starting install..."
cd
apt-get update
apt-get install git
apt-get install nginx
mkdir gopath
mkdir gopath/src
curl https://storage.googleapis.com/golang/go1.2.2.linux-amd64.tar.gz > go.tar.gz
tar -xvf go.tar.gz
export GOROOT=$HOME/go
export PATH=$PATH:$GOROOT/bin
go get github.com/go-sql-driver/mysql
cd links-as-a-service-redirectserver
git checkout GoogleAppEngine
cd 
cp links-as-a-service-redirectserver/Laas/ gopath/src -r
cp links-as-a-service-redirectserver/user_agent gopath/src -r
cd links-as-a-service-redirectserver
go build
