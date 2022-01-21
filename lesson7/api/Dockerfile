FROM golang:1.17-buster AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /api


FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /api /api
EXPOSE 8080 8081 8082 8083 8084 8085 8086 8087 8088 8089
ENTRYPOINT ["/api"]
