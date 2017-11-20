FROM haproxy:1.7

ENV HAPROXY_CONFIG_PATH /etc/haproxy/haproxy.cfg
ENV HAPROXY_PID_PATH /run/haproxy.pid

COPY go-kube-haproxy /
COPY haproxy.sh /
COPY kubeapi.sh /

ENTRYPOINT []
CMD []
