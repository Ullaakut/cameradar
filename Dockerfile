FROM golang:alpine
WORKDIR /go/src/github.com/EtixLabs/cameradar

COPY . /go/src/github.com/EtixLabs/cameradar/cameraccess

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

RUN go install

ENTRYPOINT /go/bin/cameraccess
