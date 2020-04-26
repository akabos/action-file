FROM golang:1.14-alpine3.11 AS build

RUN apk add --no-cache git

ADD . /go/action-file
WORKDIR /go/action-file

RUN go mod download
RUN go build

FROM alpine:3.11

RUN apk add --no-cache dumb-init
COPY --from=build /go/action-file/action-file /usr/local/bin

ENTRYPOINT ["/usr/bin/dumb-init", "--"]
CMD ["action-file"]
