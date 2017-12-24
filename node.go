package main

import (
  "fmt"
  apiv1 "k8s.io/api/core/v1"
)


func (tmpl *TemplateMap) UpdateNode(node *apiv1.Node) {
  for _, condition := range node.Status.Conditions {

    switch condition.Type {
    case "Ready":

      if (condition.Status != apiv1.ConditionTrue) {
        tmpl.DeleteNode(node)

        return
      }
    }
  }

  m, exists := tmpl.Nodes[node.Name]

  if !exists {
    m = &NodeMap{}
    tmpl.Nodes[node.Name] = m
  }

  for _, address := range node.Status.Addresses {

    switch address.Type {
    case "InternalIP":

      v := address.Address

      if m.Address != v {
        m.Address = v

        fmt.Printf("Update node: %s\n", node.Name)
        tmpl.Updated = true
      }

      return
    }
  }
}


func (tmpl *TemplateMap) DeleteNode(node *apiv1.Node) {
  _, exists := tmpl.Nodes[node.Name]

  if exists {
    delete(tmpl.Nodes, node.Name)

    fmt.Printf("Delete node: %s\n", node.Name)
    tmpl.Updated = true
  }
}
