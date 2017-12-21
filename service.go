package main

import (
  "fmt"
  apiv1 "k8s.io/api/core/v1"
)


func (tmpl *TemplateMap) UpdateService(service *apiv1.Service) {
  for _, value := range service.Spec.Ports {

    switch value.Protocol {
    case "TCP":

      k := ServiceKey{service.Name, value.Name}
      v := PortMap{value.NodePort, value.TargetPort.IntVal}

      if (tmpl.Services[k] != v) {
        tmpl.Services[k] = v

        fmt.Printf("Update service port: %s %d->%d\n", service.Name, value.NodePort, value.TargetPort.IntVal)
        tmpl.Updated = true
      }
    }
  }
}


func (tmpl *TemplateMap) DeleteService(service *apiv1.Service) {
  for _, value := range service.Spec.Ports {

    k := ServiceKey{service.Name, value.Name}

    _, exists := tmpl.Services[k]
    if exists {
      delete(tmpl.Services, k)

      fmt.Printf("Delete service: %s\n", service.Name)
      tmpl.Updated = true
    }
  }
}
