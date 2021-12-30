FROM golang:1.15 as modules

ADD acl/go.mod acl/go.sum /m/
RUN cd /m && go mod download

FROM golang:1.15 as builder

COPY --from=modules /go/pkg /go/pkg

RUN mkdir -p /src
ADD acl/ /src
WORKDIR /src

RUN useradd -u 10001 myapp

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
    go build -o /myapp ./

FROM busybox


COPY --from=builder /etc/passwd /etc/passwd
USER myapp

COPY --from=builder /myapp /acl
COPY --from=builder /etc/ssl/certs/ /etc/ssl/certs/
