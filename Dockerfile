FROM golang:1.8
WORKDIR /go/src/github.com/EtixLabs/cameradar

COPY . /go/src/github.com/EtixLabs/cameradar

RUN echo "deb http://ftp.debian.org/debian jessie-backports main" >> /etc/apt/sources.list \
 && apt-get update \
 && apt-get install -y libcurl4-openssl-dev nmap \
 && rm -rf /var/lib/apt/lists/* \
 && rm -rf /var/cache/apk/*

RUN go get github.com/andelf/go-curl
RUN go get github.com/pkg/errors
RUN go get gopkg.in/go-playground/validator.v9

RUN cd cameraccess ; go install -race

ENTRYPOINT /go/bin/cameraccess
