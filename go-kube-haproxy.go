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
  kubeconfigFile = flag.String("kubeconfig", "", "Kubeconfig file path")
  templateFile = flag.String("template", "", "Haproxy Go template file path")
  outFile = flag.String("output", "", "Output file path")
  pidFile = flag.String("pid", "", "PID file path")
  annotationPrefix = flag.String("prefix", "kube-haproxy.", "Annotation prefix")
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
  t := &TemplateMap {
    Services: make(map[string](*ServiceMap)),
    Nodes:    make(map[string](*NodeMap)),
    updated:  make(chan struct{}, 1),
  }


  nodeWatchlist := cache.NewListWatchFromClient(clientset.Core().RESTClient(),
    "nodes", apiv1.NamespaceAll, fields.Everything())

  _, nodeController := cache.NewInformer(
    nodeWatchlist,
    &apiv1.Node{},
    time.Second * 0,
    cache.ResourceEventHandlerFuncs{

      AddFunc: func(obj interface{}) {
        t.mux.Lock()
        t.UpdateAddresses(obj.(*apiv1.Node))
        t.NodeAnnotations(obj.(*apiv1.Node))
        t.mux.Unlock()
      },

      UpdateFunc: func(_, obj interface{}) {
        t.mux.Lock()
        t.UpdateAddresses(obj.(*apiv1.Node))
        t.NodeAnnotations(obj.(*apiv1.Node))
        t.mux.Unlock()
      },

      DeleteFunc: func(obj interface{}) {
        t.mux.Lock()
        t.DeleteNode(obj.(*apiv1.Node))
        t.mux.Unlock()
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
        t.mux.Lock()
        t.UpdatePorts(obj.(*apiv1.Service))
        t.ServiceAnnotations(obj.(*apiv1.Service))
        t.mux.Unlock()
      },

      UpdateFunc: func(_, obj interface{}) {
        t.mux.Lock()
        t.UpdatePorts(obj.(*apiv1.Service))
        t.ServiceAnnotations(obj.(*apiv1.Service))
        t.mux.Unlock()
      },

      DeleteFunc: func(obj interface{}) {
        t.mux.Lock()
        t.DeleteService(obj.(*apiv1.Service))
        t.mux.Unlock()
      },
    },
  )


  stop := make(chan struct{})

  go nodeController.Run(stop)
  go serviceController.Run(stop)

  for {
    // don't reload config more than once per 5 sec
    time.Sleep(time.Second * 5)

    select {
    case <- t.updated:
      t.mux.Lock()
      // defer t.mux.Unlock()
      err := t.WriteHaproxyConfig(configTmpl)

      if err != nil {
        fmt.Println(err)
      } else {
        callHaproxyReload(*pidFile)
      }
      t.mux.Unlock()

    default:
    }
  }
}
