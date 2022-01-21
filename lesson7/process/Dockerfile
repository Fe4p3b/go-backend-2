FROM golang:1.17-buster AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /process


FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /process /process
EXPOSE 8080
ENTRYPOINT ["/process"]
