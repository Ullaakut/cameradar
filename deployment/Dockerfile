FROM ubuntu:16.04

MAINTAINER brendan.leglaunec@etixgroup.com

ENV LD_LIBRARY_PATH="/cameradar/libraries"

RUN apt-get update && apt-get install -y \
    nmap \
    ffmpeg \
    libboost-all-dev \
    libgstreamer1.0-dev \
    gstreamer1.0-plugins-base \
    gstreamer1.0-plugins-good \
    libcurl4-openssl-dev \
    libmysqlclient20 \
    mysql-client

ADD cameradar_*_Release_Linux.tar.gz /
RUN mv cameradar_*_Release_Linux cameradar

COPY conf /cameradar/conf

COPY docker-entrypoint.sh /usr/local/bin/
RUN ln -s /usr/local/bin/docker-entrypoint.sh /entrypoint.sh
ENTRYPOINT ["docker-entrypoint.sh"]

CMD ["/cameradar/bin/cameradar", "-c", "/cameradar/conf/cameradar.conf.json"]
