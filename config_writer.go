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
func (tmpl *TemplateMap) WriteHaproxyConfig(configTmpl *template.Template) bool {
  if (tmpl.Updated) {
    tmpl.Updated = false

    err := configTmpl.Execute(io.Writer(&templateBuffer), tmpl)
    if err != nil {
      fmt.Println(err)

    } else {

      for k, v := range tmpl.Nodes {
        fmt.Printf("config node addresses: %s -> %s\n", k, v.Address)
      }

      for k, v := range tmpl.Services {
        for n, p := range v.Ports {
          fmt.Printf("config service ports: %s %s -> %d\n", k, n, p)
        }
      }

      f, err := os.OpenFile(*outFile, os.O_CREATE|os.O_RDWR, 0644)
      defer f.Close()

      if err != nil {
        fmt.Println(err)

      } else {
        written, err := templateBuffer.WriteTo(f)
        f.Truncate(written)

        if err != nil {
          fmt.Println(err)

        } else {
          return true
        }
      }
    }
  }

  return false
}
