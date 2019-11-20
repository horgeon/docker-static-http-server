FROM golang:alpine as build-env

WORKDIR /go/src/server
ADD . /go/src/server

RUN go get -d -v ./...

RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/server

FROM centurylink/ca-certs
COPY --from=build-env /go/bin/server /
CMD ["/server"]
