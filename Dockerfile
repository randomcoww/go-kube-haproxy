FROM debian:sid

COPY go-kube-haproxy /
COPY entrypoint.sh /
ENTRYPOINT ["/entrypoint.sh"]
