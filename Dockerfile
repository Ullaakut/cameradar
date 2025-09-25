# Build stage
FROM golang:alpine AS build-env

COPY . /go/src/github.com/Ullaakut/cameradar
WORKDIR /go/src/github.com/Ullaakut/cameradar/cmd/cameradar

# Necessary to install curl v7.64.0-r5.
# Fix for https://github.com/Ullaakut/cameradar/issues/247
# And a bug when curl doesn't change request method
RUN echo 'http://dl-cdn.alpinelinux.org/alpine/v3.9/main' >> /etc/apk/repositories
RUN apk update && \
    apk upgrade && \
    apk add nmap nmap-nselibs nmap-scripts \
    curl curl-dev==7.64.0-r5 \
    gcc \
    libc-dev \
    git \
    pkgconfig
ENV GO111MODULE=on
RUN go version
RUN go build -o cameradar

# Final stage
FROM alpine

RUN echo 'http://dl-cdn.alpinelinux.org/alpine/v3.9/main' >> /etc/apk/repositories

RUN apk --update add --no-cache nmap \
    nmap-nselibs \
    nmap-scripts \
    curl-dev==7.64.0-r5

WORKDIR /app/cameradar
COPY --from=build-env /go/src/github.com/Ullaakut/cameradar/dictionaries/ /app/dictionaries/
COPY --from=build-env /go/src/github.com/Ullaakut/cameradar/cmd/cameradar/ /app/cameradar/

ENV CAMERADAR_CUSTOM_ROUTES="/app/dictionaries/routes"
ENV CAMERADAR_CUSTOM_CREDENTIALS="/app/dictionaries/credentials.json"

ENTRYPOINT ["/app/cameradar/cameradar"]
