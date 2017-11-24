package main

import (
  "fmt"
  "flag"
  "time"
  "text/template"
  "os"
  // "os/exec"
  "strconv"
  "syscall"
  "io/ioutil"
  "bytes"

  apiv1 "k8s.io/api/core/v1"
  "k8s.io/client-go/kubernetes"
  "k8s.io/client-go/tools/clientcmd"
  "k8s.io/client-go/tools/cache"
  "k8s.io/apimachinery/pkg/fields"
)

type Key struct {
  ServiceName, PortName string
}

type PortMap struct {
  NodePort, Port int32
}

var (
  kubeconfigFile = flag.String("kubeconfig", "", "kubeconfig file path")
  templateFile = flag.String("template", "", "go template file path")
  outFile = flag.String("output", "", "output file path")
  pidFile = flag.String("pid", "", "pid file path")
)


// func haproxyCommand(cmd string, result *bytes.Buffer) {
//   c, err := net.Dial("unix", *socketPath)
//   defer c.Close()
//
//   if err != nil {
//     panic(err.Error())
//   }
//
//   _, err = c.Write([]byte(cmd + "\n"))
//   if err != nil {
//     panic(err.Error())
//   }
//
//   io.Copy(result, c)
// }


func main() {
  flag.Parse()

  servicesMap := make(map[Key]PortMap)
  tmpl := template.New("template")
  updated := false


  updateTemplate := func() {
    f, _ := os.OpenFile(*outFile, os.O_CREATE|os.O_WRONLY, 0644)
    defer f.Close()

    tmpl, _ = tmpl.ParseFiles(*templateFile)
    tmpl.Execute(f, servicesMap)
    fmt.Println("Update template")
  }


  callReload := func() {
    pid, err := ioutil.ReadFile(*pidFile)

    if err == nil {
      pid, err := strconv.Atoi(string(bytes.TrimSpace(pid)))

      if err == nil {
        syscall.Kill(pid, syscall.SIGUSR2)
        fmt.Println("Send kill", pid)
      }
    }
  }


  config, err := clientcmd.BuildConfigFromFlags("", *kubeconfigFile)
  if err != nil {
    panic(err.Error())
  }

  clientset, err := kubernetes.NewForConfig(config)
  if err != nil {
    panic(err.Error())
  }

  watchlist := cache.NewListWatchFromClient(clientset.Core().RESTClient(),
    "services", apiv1.NamespaceDefault, fields.Everything())

  _, controller := cache.NewInformer(
    watchlist,
    &apiv1.Service{},
    time.Second * 0,
    cache.ResourceEventHandlerFuncs{

      AddFunc: func(obj interface{}) {
        service := obj.(*apiv1.Service)

        for _, value := range service.Spec.Ports {
          servicesMap[Key{service.Name, value.Name}] = PortMap{value.NodePort, value.Port}
          fmt.Printf("Add service port: %s %d->%d\n", service.Name, value.NodePort, value.Port)
        }

        updated = true
      },

      DeleteFunc: func(obj interface{}) {
        service := obj.(*apiv1.Service)

        for _, value := range service.Spec.Ports {
          delete(servicesMap, Key{service.Name, value.Name})
        }

        fmt.Printf("Delete service: %s\n", service.Name)
        updated = true
      },

      UpdateFunc: func(_, obj interface{}) {
        service := obj.(*apiv1.Service)

        for _, value := range service.Spec.Ports {
          if (servicesMap[Key{service.Name, value.Name}] != PortMap{value.NodePort, value.Port}) {

            servicesMap[Key{service.Name, value.Name}] = PortMap{value.NodePort, value.Port}
            fmt.Printf("Update service port: %s %d->%d\n", service.Name, value.NodePort, value.Port)

            updated = true
          }
        }
      },
    },
  )

  stop := make(chan struct{})

  go controller.Run(stop)

  for {
    time.Sleep(time.Second * 5)

    if (updated) {
      updated = false

      updateTemplate()
      callReload()
    }
  }
}
