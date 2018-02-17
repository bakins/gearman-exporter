FROM alpine:3.7

ADD gearman-exporter.linux.amd64 /usr/bin/gearman-exporter
RUN chmod a+x /usr/bin/gearman-exporter

ENTRYPOINT [ "/usr/bin/gearman-exporter" ]
