package main

import (
  "fmt"
  "os"
  "text/template"
  "bytes"
  "io"
)

var (
  templateBuffer bytes.Buffer
)

// Write Template to haproxy config
func (t *TemplateMap) WriteHaproxyConfig(configTmpl *template.Template) error {

  err := configTmpl.Execute(io.Writer(&templateBuffer), t)
  if err != nil {
    return err

  } else {

    for k, v := range t.Nodes {
      fmt.Printf("config node addresses: %s -> %s\n", k, v.Address)
      for n, p := range v.Annotations {
        fmt.Printf("config node annotations: %s %s -> %s\n", k, n, p)
      }
    }

    for k, v := range t.Services {
      for n, p := range v.Ports {
        fmt.Printf("config service ports: %s %s -> %d\n", k, n, p)
      }
      for n, p := range v.Annotations {
        fmt.Printf("config service annotations: %s %s -> %s\n", k, n, p)
      }
    }

    f, err := os.OpenFile(*outFile, os.O_CREATE|os.O_RDWR, 0644)
    defer f.Close()

    if err != nil {
      return err

    } else {
      written, err := templateBuffer.WriteTo(f)
      f.Truncate(written)

      if err != nil {
        return err

      } else {
        return nil
      }
    }
  }
}
