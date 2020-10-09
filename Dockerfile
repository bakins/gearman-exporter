FROM golang:1 as builder

WORKDIR /build
COPY . .
RUN go build -mod=vendor -o /gearman-exporter cmd/gearman-exporter/main.go

FROM debian:bullseye-slim

COPY --from=builder /gearman-exporter /usr/bin/gearman-exporter
RUN chmod +x /usr/bin/gearman-exporter

CMD [ "/usr/bin/gearman-exporter" ]
