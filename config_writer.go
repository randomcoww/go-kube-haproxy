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
      fmt.Println("write nodes", tmpl.Nodes)
      fmt.Println("write services", tmpl.Services)

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
