### Haproxy side car for Kubernetes node port load balancing

A mapping of service target ports and node ports are collected for use in building Haproxy configuration from a supplied Go template.

When one or more mapped values change, a new configuration is written and seamless reload is called on Haproxy. Haproxy 1.8 or newer is required. Shared PID space with HAProxy is required. For Kubernetes, the --docker-disable-shared-pid=false option is currently needed on Kubelet to allow this.

### Arguments
    -kubeconfig Path to kubeconfig file
    -template Path to Go template file for Haproxy
    -outfile Haproxy config path
    -pid Path to Haproxy PID file
    -prefix Read and add node and service annotations starting with this prefix into mapping for general use

### Examples

Make services on node ports available over Haproxy via target port. Add all node members by node port to backend.

    {{range $name, $s := $.Services}}{{range $portname, $p := $s.Ports}}frontend {{$name}}_{{$portname}}
      default_backend {{$name}}_{{$portname}}
      bind *:{{$p.TargetPort}}
      maxconn 2000
    backend {{$name}}_{{$portname}}
      {{range $nodename, $n := $.Nodes}}server {{$nodename}} {{$n.Address}}:{{$p.NodePort}} check
      {{end}}
    {{end}}{{end}}

Add path begins with reverse proxy rules using service annotations

    frontend http-in
      bind 0.0.0.0:80
      {{range $name, $s := $.Services}}{{range $k, $annotation := $s.Annotations}}acl acl_{{$name}}_{{$k}} {{$annotation}}
      use_backend {{$name}}_{{$k}} if acl_{{$name}}_{{$k}}
      {{end}}{{end}}
