package main

import (
  "fmt"
  "os"
  "text/template"
  "path/filepath"
)

// Write Template to haproxy config
func (tmpl *TemplateMap) WriteHaproxyConfig(templateFile string) {
  f, _ := os.OpenFile(*outFile, os.O_CREATE|os.O_RDWR, 0644)
  defer f.Close()
  f.Truncate(0)

  _, file := filepath.Split(templateFile)

  t, err := template.New(file).ParseFiles(templateFile)
  if err != nil {
    fmt.Println(err)
  }

  err = t.Execute(f, tmpl)
  if err != nil {
    fmt.Println(err)

  } else {
    fmt.Println("write nodes", tmpl.Nodes)
    fmt.Println("write services", tmpl.Services)
  }
}
