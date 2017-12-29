package main

import (
  "strings"
)


type Annotations interface{
  GetAnnotations() map[string]string
  SetAnnotations(a map[string]string)
}


func (s *ServiceMap) GetAnnotations() map[string]string {
  return s.Annotations
}

func (s *ServiceMap) SetAnnotations(a map[string]string) {
  s.Annotations = a
}

func (n *NodeMap) GetAnnotations() map[string]string {
  return n.Annotations
}

func (n *NodeMap) SetAnnotations(a map[string]string) {
  n.Annotations = a
}


// annotations
func (t *TemplateMap) UpdateAnnotations(a map[string]string, m Annotations, isNew bool) {

  newAnnotations := make(map[string]string)
  updated := false

  for k, v := range a {

    if strings.HasPrefix(k, *annotationPrefix) {
      k = strings.TrimLeft(k, *annotationPrefix)

      if k == "" {
        continue
      }

      if !isNew && m.GetAnnotations()[k] != v {
        updated = true
      }

      newAnnotations[k] = v
    }
  }

  if isNew || updated || len(newAnnotations) != len(m.GetAnnotations()) {
    m.SetAnnotations(newAnnotations)
    t.setUpdated()
  }
}
