package main

import (
)

type TemplateMap struct {
  Services  map[string](*ServiceMap)
  Nodes     map[string](*NodeMap)
  Updated   bool
}


type ServiceMap struct {
  Ports       map[string]PortMap
  Annotations map[string]string
}

type NodeMap struct {
  Address     string
  Annotations map[string]string
}

type PortMap struct {
  NodePort    int32
  TargetPort  int32
}
