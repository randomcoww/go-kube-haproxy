package main

import (
  "fmt"
  "os"
  "text/template"
)

// Write Template to haproxy config
func (tmpl *TemplateMap) WriteHaproxyConfig(configTmpl *template.Template) {
  f, _ := os.OpenFile(*outFile, os.O_CREATE|os.O_RDWR, 0644)
  defer f.Close()
  f.Truncate(0)

  err := configTmpl.Execute(f, tmpl)
  if err != nil {
    fmt.Println(err)

  } else {
    fmt.Println("write nodes", tmpl.Nodes)
    fmt.Println("write services", tmpl.Services)
  }
}
