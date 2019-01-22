# Build stage
FROM golang:alpine AS build-env

COPY . /go/src/github.com/Ullaakut/cameradar
WORKDIR /go/src/github.com/Ullaakut/cameradar/cameradar

RUN apk update && \
    apk upgrade && \
    apk add nmap nmap-nselibs nmap-scripts \
    curl curl-dev \
    gcc \
    libc-dev \
    git \
    pkgconfig
ENV DEP_VERSION="0.5.0"
RUN curl -L -s https://github.com/golang/dep/releases/download/v${DEP_VERSION}/dep-linux-amd64 -o $GOPATH/bin/dep
RUN chmod +x $GOPATH/bin/dep
RUN dep ensure
RUN go build -o cameradar

# Final stage
FROM alpine

RUN apk --update add --no-cache nmap \
    nmap-nselibs \
    nmap-scripts \
    curl-dev

WORKDIR /app/cameradar
COPY --from=build-env /go/src/github.com/Ullaakut/cameradar/dictionaries/ /app/dictionaries/
COPY --from=build-env /go/src/github.com/Ullaakut/cameradar/cameradar/ /app/cameradar/
ENTRYPOINT ["/app/cameradar/cameradar", "-r", "/app/dictionaries/routes", "-c", "/app/dictionaries/credentials.json"]
