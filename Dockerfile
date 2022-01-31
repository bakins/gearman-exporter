FROM golang:1-alpine as build
WORKDIR /go/gearman-exporter
COPY . .
RUN go build -o /gearman-exporter ./cmd/gearman-exporter

FROM alpine:3
COPY --from=build /gearman-exporter /usr/bin/gearman-exporter
ENTRYPOINT [ "/usr/bin/gearman-exporter" ]
