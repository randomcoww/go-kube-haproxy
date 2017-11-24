FROM haproxy:rc

COPY go-kube-haproxy /
COPY kubeapi.sh /
