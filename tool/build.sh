#!/bin/bash

cd $GOPATH/src/github.com/mathewreny/sftp/
go build -o gentool ./tool/gen.go && 
cat ./tool/requests.gen |
./gentool |
gofmt /dev/stdin > ./requests.go

rm -f gentool
