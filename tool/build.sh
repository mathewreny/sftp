#!/bin/bash

cd $GOPATH/src/github.com/mathewreny/sftp/
go build -o gentool ./tool/gen.go && 
	cat ./tool/requests.gen |
	./gentool |
	gofmt /dev/stdin > ./requests.go

rm -f gentool

go build -o utilgentool ./tool/utilgen.go &&
	cat ./tool/requests.gen |
	./utilgentool |
	gofmt /dev/stdin > ./packets.go

rm -f utilgentool
