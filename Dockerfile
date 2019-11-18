FROM golang:1.12 as builder
WORKDIR /go/src/github.com/camptocamp/prometheus-kerberos-exporter
COPY . .
RUN make prometheus-kerberos-exporter

FROM scratch
COPY --from=builder /go/src/github.com/camptocamp/prometheus-kerberos-exporter/prometheus-kerberos-exporter /
ENTRYPOINT ["/prometheus-kerberos-exporter"]
CMD [""]
