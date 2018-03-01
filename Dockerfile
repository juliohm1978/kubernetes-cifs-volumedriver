FROM busybox:1.28.1

COPY juliohm~cifs /
COPY install.sh /usr/local/bin/
RUN chmod +x /usr/local/bin/install.sh

CMD ["/usr/local/bin/install.sh"]
