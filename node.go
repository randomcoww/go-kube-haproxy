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

  for _, address := range node.Status.Addresses {

    k := NodeKey{node.Name}
    v := IPMap{address.Address}

    switch address.Type {
    case "InternalIP":

      if (tmpl.Nodes[k] != v) {
        tmpl.Nodes[k] = v

        fmt.Printf("Update node: %s->%s\n", node.Name, address.Address)
        tmpl.Updated = true
      }
    }
  }
}


func (tmpl *TemplateMap) DeleteNode(node *apiv1.Node) {
  k := NodeKey{node.Name}

  _, exists := tmpl.Nodes[k]
  if exists {
    delete(tmpl.Nodes, k)

    fmt.Printf("Delete node: %s\n", node.Name)
    tmpl.Updated = true
  }
}
