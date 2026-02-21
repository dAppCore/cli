# Host UK — Core CLI Container
# Multi-stage build: Go binary in distroless-style Alpine
#
# Build: docker build -f docker/Dockerfile.core -t lthn/core:latest .

FROM golang:1.26-alpine AS build

RUN apk add --no-cache git ca-certificates

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -trimpath -ldflags '-s -w' -o /core .

FROM alpine:3.21
RUN apk add --no-cache ca-certificates
COPY --from=build /core /usr/local/bin/core
RUN adduser -D -h /home/core core
USER core
ENTRYPOINT ["core"]
