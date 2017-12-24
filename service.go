package main

import (
  "fmt"
  apiv1 "k8s.io/api/core/v1"
)


func (tmpl *TemplateMap) UpdateService(service *apiv1.Service) {
  m, exists := tmpl.Services[service.Name]

  if !exists {
    m = &ServiceMap{}
    tmpl.Services[service.Name] = m
  }

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

      if exists && m.Ports[k] != v {
        updated = true
      }

      newPorts[k] = v
    }
  }

  if !exists || updated || len(newPorts) != len(m.Ports) {
    m.Ports = newPorts

    fmt.Printf("Update service: %s\n", service.Name)
    tmpl.Updated = true
  }
}


func (tmpl *TemplateMap) DeleteService(service *apiv1.Service) {
  _, exists := tmpl.Services[service.Name]

  if exists {
    delete(tmpl.Services, service.Name)

    fmt.Printf("Delete service: %s\n", service.Name)
    tmpl.Updated = true
  }
}
