#!/bin/bash

cd $GOPATH/src/vendor/golang.org/x/lint/golint
go build
mv golint $GOPATH/bin

cd $GOPATH/src/vendor/golang.org/x/tools/cmd/goimports
go build
mv goimports $GOPATH/bin

cd $GOPATH/src/vendor/golang.org/x/tools/cmd/gorename
go build
mv gorename $GOPATH/bin

cd $GOPATH/src/vendor/golang.org/x/tools/cmd/guru
go build
mv guru $GOPATH/bin

cd $GOPATH/src/vendor/github.com/nsf/gocode
go build
mv gocode $GOPATH/bin

cd $GOPATH/src/vendor/github.com/rogpeppe/godef
go build
mv godef $GOPATH/bin

cd $GOPATH/src/vendor/github.com/lukehoban/go-outline
go build
mv go-outline $GOPATH/bin

cd $GOPATH/src/vendor/github.com/sqs/goreturns
go build
mv goreturns $GOPATH/bin

cd $GOPATH/src/vendor/github.com/tpng/gopkgs
go build
mv gopkgs $GOPATH/bin

cd $GOPATH/src/vendor/github.com/newhook/go-symbols
go build
mv go-symbols $GOPATH/bin

# cobra tool
cd $GOPATH/src/vendor/github.com/spf13/cobra/cobra
go build
mv cobra $GOPATH/bin
