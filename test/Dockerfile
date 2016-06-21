FROM ubuntu:15.10

MAINTAINER brendan.leglaunec@etixgroup.com

ENV LD_LIBRARY_PATH="/cctv/libraries"

# install go
RUN apt-get update && apt-get install -y make git wget curl
RUN wget https://storage.googleapis.com/golang/go1.6.linux-amd64.tar.gz
RUN tar -C /usr/local -xzf go1.6.linux-amd64.tar.gz
# set variable env
ENV GOPATH=/go
ENV PATH=$PATH:/go/bin
ENV PATH=$PATH:/usr/local/go/bin
ENV LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/usr/local/lib

# needed for cameradar
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

ADD cctv_*_Debug_Linux.tar.gz /
RUN mv cctv_*_Debug_Linux cctv

# create cameradaratest folder in go src path
RUN mkdir -p /go/src/cameradartest
ADD ./conf /conf
ADD ./docker/run_cameradartest.sh /run.sh

# get go deps
RUN go get github.com/go-sql-driver/mysql

RUN mkdir /thumbnails
WORKDIR /go/src/cameradartest
CMD ["/run.sh"]
