FROM scratch

ADD gearman-exporter.linux.amd64 /gearman-exporter

ENTRYPOINT [ "/gearman-exporter" ]
