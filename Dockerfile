FROM golang:alpine AS builder
LABEL stage=gobuilder

ENV CGO_ENABLED=0
ENV GOOS=linux

RUN apk update --no-cache && apk add --no-cache tzdata

WORKDIR /build

ADD go.mod go.sum .
RUN go mod download

COPY . .
RUN go build -ldflags="-s -w" -o /app/shortener ./cmd/main.go

FROM alpine
RUN apk update --no-cache && apk add --no-cache ca-certificates

ENV TZ=UTC
RUN ln -sf /usr/share/zoneinfo/UTC /etc/localtime

WORKDIR /app
COPY --from=builder /app/shortener /app/shortener

EXPOSE 8080

ENTRYPOINT ["./shortener"]
