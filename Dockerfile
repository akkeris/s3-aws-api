FROM golang:1.8-alpine
RUN apk update
RUN apk add openssl ca-certificates git
RUN mkdir -p /go/src/s3-aws-api
ADD server.go  /go/src/s3-aws-api/server.go
ADD create.sql /go/src/s3-aws-api/create.sql
ADD iam /go/src/s3-aws-api/iam
ADD db /go/src/s3-aws-api/db
ADD structs /go/src/s3-aws-api/structs
ADD s3 /go/src/s3-aws-api/s3
ADD utils /go/src/s3-aws-api/utils
ADD build.sh /build.sh
RUN chmod +x /build.sh
RUN /build.sh
CMD ["/go/src/s3-aws-api/server"]
EXPOSE 3500
