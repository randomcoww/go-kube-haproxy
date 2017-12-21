package main

type ServiceKey struct {
  ServiceName string
  PortName    string
}

type PortMap struct {
  NodePort  int32
  Port      int32
}


type NodeKey struct {
  NodeName  string
}

type IPMap struct {
  InternalIP string
}


type TemplateMap struct {
  Services  map[ServiceKey]PortMap
  Nodes     map[NodeKey]IPMap
  Updated   bool
}
