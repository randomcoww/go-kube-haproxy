package main

import (
  "fmt"
  "strings"
  apiv1 "k8s.io/api/core/v1"
)


func (t *TemplateMap) nodeMap(node *apiv1.Node) (*NodeMap, bool) {
  m, exists := t.Nodes[node.Name]

  if !exists {
    m = &NodeMap{}
    t.Nodes[node.Name] = m
  }

  return m, !exists
}

// node addresses
func (t *TemplateMap) UpdateAddresses(node *apiv1.Node) {
  for _, condition := range node.Status.Conditions {

    switch condition.Type {
    case "Ready":

      if (condition.Status != apiv1.ConditionTrue) {
        t.DeleteNode(node)

        return
      }
    }
  }

  m, _ := t.nodeMap(node)

  for _, address := range node.Status.Addresses {

    switch address.Type {
    case "InternalIP":

      v := address.Address

      if m.Address != v {
        m.Address = v

        fmt.Printf("Update node addresses: %s\n", node.Name)
        t.Updated = true
      }

      return
    }
  }
}

// service annotations
func (t *TemplateMap) NodeAnnotations(node *apiv1.Node) {
  m, new := t.nodeMap(node)

  newAnnotations := make(map[string]string)
  updated := false

  for k, v := range node.Annotations {

    if strings.HasPrefix(k, "kube_haproxy.") {
      k = strings.TrimLeft(k, "kube_haproxy.")

      if k == "" {
        continue
      }

      if !new && m.Annotations[k] != v {
        updated = true
      }

      newAnnotations[k] = v
    }
  }

  if new || updated || len(newAnnotations) != len(m.Annotations) {
    m.Annotations = newAnnotations

    fmt.Printf("Update node annotations: %s\n", node.Name)
    t.Updated = true
  }
}


func (t *TemplateMap) DeleteNode(node *apiv1.Node) {
  _, exists := t.Nodes[node.Name]

  if exists {
    delete(t.Nodes, node.Name)

    fmt.Printf("Delete node: %s\n", node.Name)
    t.Updated = true
  }
}
