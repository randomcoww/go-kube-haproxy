FROM haproxy:rc

ENV HAPROXY_CONFIG_PATH /etc/haproxy/haproxy.cfg

COPY go-kube-haproxy /
COPY kubeapi.sh /
