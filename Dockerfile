FROM golang:latest as BUILDER

# build binary
COPY . /go/src/github.com/opensourceways/software-package-server
RUN cd /go/src/github.com/opensourceways/software-package-server && GO111MODULE=on CGO_ENABLED=0 go build

# copy binary config and utils
FROM alpine:latest
WORKDIR /opt/app/

COPY  --from=BUILDER /go/src/github.com/opensourceways/software-package-server/software-package-server /opt/app

ENTRYPOINT ["/opt/app/software-package-server"]
