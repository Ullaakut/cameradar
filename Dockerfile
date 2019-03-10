# Build stage
FROM golang:alpine AS build-env

COPY . /go/src/github.com/ullaakut/cameradar
WORKDIR /go/src/github.com/ullaakut/cameradar/cameradar

RUN apk update && \
    apk upgrade && \
    apk add nmap nmap-nselibs nmap-scripts \
    curl curl-dev \
    gcc \
    libc-dev \
    git \
    pkgconfig
ENV GO111MODULE=on
RUN go version
RUN go build -o cameradar

# Final stage
FROM alpine

RUN apk --update add --no-cache nmap \
    nmap-nselibs \
    nmap-scripts \
    curl-dev

WORKDIR /app/cameradar
COPY --from=build-env /go/src/github.com/ullaakut/cameradar/dictionaries/ /app/dictionaries/
COPY --from=build-env /go/src/github.com/ullaakut/cameradar/cameradar/ /app/cameradar/

ENV CAMERADAR_CUSTOM_ROUTES="/app/dictionaries/routes"
ENV CAMERADAR_CUSTOM_CREDENTIALS="/app/dictionaries/credentials.json"

ENTRYPOINT ["/app/cameradar/cameradar"]
