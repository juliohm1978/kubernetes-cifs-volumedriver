FROM golang:alpine AS build-env
RUN apk --no-cache add build-base git gcc

ADD . /kubernetes-cifs-volumedriver
WORKDIR /kubernetes-cifs-volumedriver
RUN go build -a -installsuffix cgo

FROM busybox:1.28.4

ENV VENDOR=juliohm
ENV DRIVER=cifs

COPY --from=build-env /kubernetes-cifs-volumedriver /
COPY install.sh /usr/local/bin/
RUN chmod +x /usr/local/bin/install.sh

CMD ["/usr/local/bin/install.sh"]
