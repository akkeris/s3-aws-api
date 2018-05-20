#!/bin/sh

cd /go/src
go get  "github.com/lib/pq"
go get  "github.com/aws/aws-sdk-go/aws"
go get  "github.com/aws/aws-sdk-go/aws/session"
go get  "github.com/aws/aws-sdk-go/service/s3"
go get  "github.com/go-martini/martini"
go get  "github.com/martini-contrib/binding"
go get  "github.com/martini-contrib/render"
go get  "github.com/nu7hatch/gouuid"
go get  "github.com/bitly/go-simplejson"
cd /go/src/s3-aws-api
go build server.go

