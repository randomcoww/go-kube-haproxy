package main

import (
  "fmt"
  "flag"
  "time"
  "text/template"
  "os"
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

// service
type ServiceKey struct {
  ServiceName, PortName string
}

type PortMap struct {
  NodePort, Port int32
}


// node
type NodeKey struct {
  NodeName string
}

type IPMap struct {
  InternalIP string
}

// template
type Template struct {
  Services map[ServiceKey]PortMap
  Nodes map[NodeKey]IPMap
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

  servicesMap := make(map[ServiceKey]PortMap)
  nodesMap := make(map[NodeKey]IPMap)

  tmpl := template.New("template")
  updated := false


  updateTemplate := func() {
    f, _ := os.OpenFile(*outFile, os.O_CREATE|os.O_WRONLY, 0644)
    defer f.Close()

    tmpl, _ = tmpl.ParseFiles(*templateFile)
    fmt.Println("Update template", Template{ Services: servicesMap, Nodes: nodesMap })

    err := tmpl.Execute(f, Template{ Services: servicesMap, Nodes: nodesMap })
    if err != nil {
      fmt.Println(err)
    }
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


  nodeWatchlist := cache.NewListWatchFromClient(clientset.Core().RESTClient(),
    "nodes", apiv1.NamespaceAll, fields.Everything())

  _, nodeController := cache.NewInformer(
    nodeWatchlist,
    &apiv1.Node{},
    time.Second * 0,
    cache.ResourceEventHandlerFuncs{

      AddFunc: func(obj interface{}) {
        node := obj.(*apiv1.Node)
        if node.Spec.Unschedulable {
					return
        }

        for _, address := range node.Status.Addresses {
          if (address.Type == "InternalIP") {

            nodesMap[NodeKey{node.Name}] = IPMap{address.Address}
            fmt.Printf("Add node: %s->%s\n", node.Name, address.Address)
            updated = true
          }
        }
      },

      DeleteFunc: func(obj interface{}) {
        node := obj.(*apiv1.Node)

        delete(nodesMap, NodeKey{node.Name})
        fmt.Printf("Delete node: %s\n", node.Name)
        updated = true
      },

      UpdateFunc: func(_, obj interface{}) {
        node := obj.(*apiv1.Node)
        if node.Spec.Unschedulable {
					return
        }

        for _, address := range node.Status.Addresses {
          if (address.Type == "InternalIP" &&
            nodesMap[NodeKey{node.Name}] != IPMap{address.Address}) {

            nodesMap[NodeKey{node.Name}] = IPMap{address.Address}
            fmt.Printf("Update node: %s->%s\n", node.Name, address.Address)
            updated = true
          }
        }
      },
    },
  )


  serviceWatchlist := cache.NewListWatchFromClient(clientset.Core().RESTClient(),
    "services", apiv1.NamespaceDefault, fields.Everything())

  _, serviceController := cache.NewInformer(
    serviceWatchlist,
    &apiv1.Service{},
    time.Second * 0,
    cache.ResourceEventHandlerFuncs{

      AddFunc: func(obj interface{}) {
        service := obj.(*apiv1.Service)

        for _, value := range service.Spec.Ports {

          servicesMap[ServiceKey{service.Name, value.Name}] = PortMap{value.NodePort, value.TargetPort.IntVal}
          fmt.Printf("Add service port: %s %d->%d\n", service.Name, value.NodePort, value.TargetPort.IntVal)
        }

        updated = true
      },

      DeleteFunc: func(obj interface{}) {
        service := obj.(*apiv1.Service)

        for _, value := range service.Spec.Ports {
          delete(servicesMap, ServiceKey{service.Name, value.Name})
        }

        fmt.Printf("Delete service: %s\n", service.Name)
        updated = true
      },

      UpdateFunc: func(_, obj interface{}) {
        service := obj.(*apiv1.Service)

        for _, value := range service.Spec.Ports {
          if (servicesMap[ServiceKey{service.Name, value.Name}] != PortMap{value.NodePort, value.TargetPort.IntVal}) {

            servicesMap[ServiceKey{service.Name, value.Name}] = PortMap{value.NodePort, value.TargetPort.IntVal}
            fmt.Printf("Update service port: %s %d->%d\n", service.Name, value.NodePort, value.TargetPort.IntVal)

            updated = true
          }
        }
      },
    },
  )

  stop := make(chan struct{})

  go nodeController.Run(stop)
  go serviceController.Run(stop)

  for {
    time.Sleep(time.Second * 5)

    if (updated) {
      updated = false

      updateTemplate()
      callReload()
    }
  }
}
