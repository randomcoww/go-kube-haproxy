FROM debian:sid

COPY go-kube-haproxy /
ENTRYPOINT ["/go-kube-haproxy"]
