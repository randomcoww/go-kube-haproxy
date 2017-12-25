package main

import (
  "fmt"
  // "strings"
  apiv1 "k8s.io/api/core/v1"
)


func (t *TemplateMap) serviceMap(service *apiv1.Service) (*ServiceMap, bool) {
  m, exists := t.Services[service.Name]

  if !exists {
    m = &ServiceMap{}
    t.Services[service.Name] = m
  }

  return m, !exists
}

// service ports
func (t *TemplateMap) UpdatePorts(service *apiv1.Service) {
  m, new := t.serviceMap(service)

  newPorts := make(map[string]PortMap)
  updated := false

  for _, port := range service.Spec.Ports {

    switch port.Protocol {
    case "TCP":

      k := port.Name
      v := PortMap{
        NodePort:   port.NodePort,
        TargetPort: port.TargetPort.IntVal,
      }

      if !new && m.Ports[k] != v {
        updated = true
      }

      newPorts[k] = v
    }
  }

  if new || updated || len(newPorts) != len(m.Ports) {
    m.Ports = newPorts

    fmt.Printf("Update service ports: %s\n", service.Name)
    t.Updated = true
  }
}

// service annotations
func (t *TemplateMap) ServiceAnnotations(service *apiv1.Service) {
  m, isNew := t.serviceMap(service)

  t.UpdateAnnotations(service.Annotations, m, isNew)
}


func (t *TemplateMap) DeleteService(service *apiv1.Service) {
  _, exists := t.Services[service.Name]

  if exists {
    delete(t.Services, service.Name)

    fmt.Printf("Delete service: %s\n", service.Name)
    t.Updated = true
  }
}
