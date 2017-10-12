# Build stage
FROM golang:alpine AS build-env

COPY . /go/src/github.com/EtixLabs/cameradar
WORKDIR /go/src/github.com/EtixLabs/cameradar/cameradar

RUN apk update && \
    apk upgrade && \
    apk add nmap nmap-nselibs nmap-scripts \
            curl curl-dev \
            gcc \
            libc-dev \
            git \
            pkgconfig

RUN curl https://glide.sh/get | sh && glide install
RUN go build -o cameradar

# Final stage
FROM alpine

RUN apk --update add --no-cache nmap nmap-nselibs nmap-scripts \
            curl-dev

WORKDIR /app/cameradar
COPY --from=build-env /go/src/github.com/EtixLabs/cameradar/ /app/
ENTRYPOINT ["/app/cameradar/cameradar"]