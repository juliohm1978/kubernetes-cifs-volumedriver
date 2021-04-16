FROM golang:1.16.3 AS build-env
ARG TARGETOS
ARG TARGETARCH
ENV GOOS=${TARGETOS}
ENV GOARCH=${TARGETARCH}

RUN apt-get update && apt-get install -y git gcc
ADD . /kubernetes-cifs-volumedriver
WORKDIR /kubernetes-cifs-volumedriver

## Running these in separate steps gives a better error 
## output indicating which one actually failed.
RUN go build -a -installsuffix cgo
RUN go test

FROM busybox:1.32.0

ENV VENDOR=juliohm
ENV DRIVER=cifs

COPY --from=build-env /kubernetes-cifs-volumedriver/kubernetes-cifs-volumedriver /
COPY install.sh /usr/local/bin/
RUN chmod +x /usr/local/bin/install.sh

CMD ["/usr/local/bin/install.sh"]
