FROM ubuntu:15.10

MAINTAINER brendan.leglaunec@etixgroup.com

ENV LD_LIBRARY_PATH="/cameradar/libraries"

# Manually install Go
RUN apt-get update && apt-get install -y make git wget curl
RUN wget https://storage.googleapis.com/golang/go1.6.linux-amd64.tar.gz
RUN tar -C /usr/local -xzf go1.6.linux-amd64.tar.gz

ENV GOPATH=/cameradartest/go
ENV PATH=$PATH:/go/bin
ENV PATH=$PATH:/usr/local/go/bin
ENV LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/usr/local/lib

RUN apt-get update && apt-get install -y \
    nmap \
    libmysqlclient18 \
    ffmpeg \
    mysql-client \
    libgstreamer1.0-dev \
    gstreamer1.0-plugins-base \
    gstreamer1.0-plugins-good \
    libcurl4-openssl-dev

RUN apt-get install -y psmisc

ADD cameradar_*_Debug_Linux.tar.gz /
RUN mv cameradar_*_Debug_Linux cameradar

RUN mkdir -p /cameradartest/go/src/cameradartest
COPY src/*.go /cameradartest/go/src/cameradartest/
COPY ./conf /conf
ADD ./docker/run_cameradartest.sh /run.sh

RUN go get github.com/go-sql-driver/mysql

RUN mkdir /thumbnails
WORKDIR /cameradartest/go/src/cameradartest
RUN go build -o cameradartest *.go
CMD ["/run.sh"]
