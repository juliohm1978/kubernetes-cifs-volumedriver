FROM busybox:1.28.4

ENV VENDOR=juliohm
ENV DRIVER=cifs

COPY kubernetes-cifs-volumedriver /
COPY install.sh /usr/local/bin/
RUN chmod +x /usr/local/bin/install.sh

CMD ["/usr/local/bin/install.sh"]
