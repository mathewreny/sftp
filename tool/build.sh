#!/bin/bash

cd $GOPATH/src/github.com/mathewreny/sftp/
go build -o clientrequestgentool ./tool/clientrequestgen.go && 
	cat ./tool/requests.gen |
	./clientrequestgentool |
	gofmt /dev/stdin > ./requests.go

rm -f clientrequestgentool

go build -o packetgentool ./tool/packetgen.go &&
	cat ./tool/packetrequests.gen ./tool/packetresponses.gen |
	./packetgentool |
	gofmt /dev/stdin > ./packets.go

rm -f packetgentool
