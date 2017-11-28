FROM alpine:latest

COPY go-kube-haproxy /
COPY entrypoint.sh /
ENTRYPOINT ["/entrypoint.sh"]
