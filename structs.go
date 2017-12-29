package main

import (
  "sync"
)

type TemplateMap struct {
  Services  map[string](*ServiceMap)
  Nodes     map[string](*NodeMap)
  mux       sync.Mutex
  updated   chan struct{}
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
