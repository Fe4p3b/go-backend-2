FROM golang:1.15 as modules

ADD go.mod go.sum /m/
RUN cd /m && go mod download

FROM golang:1.15 as builder

COPY --from=modules /go/pkg /go/pkg

RUN mkdir -p /src
ADD . /src
WORKDIR /src

RUN useradd -u 10001 myapp

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
    go build -o /myapp ./

# Готовим пробный файл статики
RUN mkdir -p /test_static && touch /test_static/index.html
RUN echo "Hello, world!" > /test_static/index.html

FROM busybox

ENV PORT 8080
ENV STATICS_PATH /test_static

COPY --from=builder /test_static /test_static

COPY --from=builder /etc/passwd /etc/passwd
USER myapp

COPY --from=builder /myapp /myapp
COPY --from=builder /etc/ssl/certs/ /etc/ssl/certs/

CMD ["/myapp"]

