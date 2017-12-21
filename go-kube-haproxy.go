package main

import (
  "fmt"
  "flag"
  "time"
  "strconv"
  "syscall"
  "bytes"
  "io/ioutil"
  "text/template"
  "path/filepath"

  apiv1 "k8s.io/api/core/v1"
  "k8s.io/client-go/kubernetes"
  "k8s.io/client-go/tools/clientcmd"
  "k8s.io/client-go/tools/cache"
  "k8s.io/apimachinery/pkg/fields"
)


var (
  kubeconfigFile = flag.String("kubeconfig", "", "kubeconfig file path")
  templateFile = flag.String("template", "", "go template file path")
  outFile = flag.String("output", "", "output file path")
  pidFile = flag.String("pid", "", "pid file path")
)


// Do Haproxy (1.8) seamless reload with USR2 signal
func callHaproxyReload(pidFile string) {
  pid, err := ioutil.ReadFile(pidFile)

  if err == nil {
    pid, err := strconv.Atoi(string(bytes.TrimSpace(pid)))

    if err == nil {
      syscall.Kill(pid, syscall.SIGUSR2)
      fmt.Println("Send kill", pid)
    }
  }
}


func main() {
  flag.Parse()

  config, err := clientcmd.BuildConfigFromFlags("", *kubeconfigFile)
  if err != nil {
    panic(err.Error())
  }

  clientset, err := kubernetes.NewForConfig(config)
  if err != nil {
    panic(err.Error())
  }

  // go template for haproxy
  _, fileName := filepath.Split(*templateFile)
  configTmpl, err := template.New(fileName).ParseFiles(*templateFile)
  if err != nil {
    panic(err.Error())
  }

  // template map from service
  tmpl := &TemplateMap {
    Services: make(map[ServiceKey]PortMap),
    Nodes:    make(map[NodeKey]IPMap),
    Updated:  false,
  }


  nodeWatchlist := cache.NewListWatchFromClient(clientset.Core().RESTClient(),
    "nodes", apiv1.NamespaceAll, fields.Everything())

  _, nodeController := cache.NewInformer(
    nodeWatchlist,
    &apiv1.Node{},
    time.Second * 0,
    cache.ResourceEventHandlerFuncs{

      AddFunc: func(obj interface{}) {
        tmpl.UpdateNode(obj.(*apiv1.Node))
      },

      UpdateFunc: func(_, obj interface{}) {
        tmpl.UpdateNode(obj.(*apiv1.Node))
      },

      DeleteFunc: func(obj interface{}) {
        tmpl.DeleteNode(obj.(*apiv1.Node))
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
        tmpl.UpdateService(obj.(*apiv1.Service))
      },

      UpdateFunc: func(_, obj interface{}) {
        tmpl.UpdateService(obj.(*apiv1.Service))
      },

      DeleteFunc: func(obj interface{}) {
        tmpl.DeleteService(obj.(*apiv1.Service))
      },
    },
  )


  stop := make(chan struct{})

  go nodeController.Run(stop)
  go serviceController.Run(stop)

  for {
    time.Sleep(time.Second * 5)

    if (tmpl.Updated) {
      tmpl.Updated = false

      tmpl.WriteHaproxyConfig(configTmpl)
      callHaproxyReload(*pidFile)
    }
  }
}
