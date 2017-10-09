# Build stage
FROM golang:alpine AS build-env

COPY . /go/src/github.com/EtixLabs/cameradar
WORKDIR /go/src/github.com/EtixLabs/cameradar/cameraccess

RUN apk update && \
    apk upgrade && \
    apk add nmap nmap-nselibs nmap-scripts \
            curl-dev \
            gcc \
            libc-dev \
            git \
            pkgconfig

RUN go get github.com/andelf/go-curl
RUN go get github.com/pkg/errors
RUN go get gopkg.in/go-playground/validator.v9
RUN go get github.com/jessevdk/go-flags
RUN go get github.com/fatih/color
RUN go get github.com/gernest/wow

RUN go build -o cameraccess

# Final stage
FROM alpine

RUN apk --update add --no-cache nmap nmap-nselibs nmap-scripts \
            curl-dev

WORKDIR /app/cameraccess
COPY --from=build-env /go/src/github.com/EtixLabs/cameradar/ /app/
ENTRYPOINT ["/app/cameraccess/cameraccess"]