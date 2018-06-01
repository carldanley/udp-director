FROM alpine:latest as alpine

RUN apk --no-cache add tzdata zip ca-certificates
RUN rm -rf /var/cache/apk/*

WORKDIR /usr/share/zoneinfo
RUN zip -r -0 /zoneinfo.zip .
ENV ZONEINFO /zoneinfo.zip

WORKDIR /
ADD udp-director /bin/
CMD ["udp-director", "--help"]
