FROM alpine

RUN apk --update add --no-cache nmap \
    nmap-nselibs \
    nmap-scripts

WORKDIR /app/cameradar

COPY cameradar /app/cameradar/cameradar

ENV CAMERADAR_CUSTOM_ROUTES="/app/dictionaries/routes"
ENV CAMERADAR_CUSTOM_CREDENTIALS="/app/dictionaries/credentials.json"

ENTRYPOINT ["/app/cameradar/cameradar"]
