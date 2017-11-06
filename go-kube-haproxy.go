package main

import (
  // "fmt"
  "flag"
  "time"

  // appsv1beta1 "k8s.io/api/apps/v1beta1"
  apiv1 "k8s.io/api/core/v1"
  "k8s.io/client-go/kubernetes"
  "k8s.io/client-go/tools/clientcmd"
  "k8s.io/client-go/tools/cache"
  "k8s.io/apimachinery/pkg/fields"

  "text/template"
  // "io"
  // "io/ioutil"
  "os"
)

type Key struct {
  ServiceName, PortName string
}

type PortMap struct {
  NodePort, Port int32
}

var (
  kubeconfigPath = flag.String("kubeconfig", "", "kubeconfig file path")
  templatePath = flag.String("template", "", "go template file path")
  outPath = flag.String("output", "", "output file path")

  servicesMap map[Key]PortMap
  output string
)


func main() {
  flag.Parse()

  config, err := clientcmd.BuildConfigFromFlags("", *kubeconfigPath)
  if err != nil {
    panic(err.Error())
  }

  clientset, err := kubernetes.NewForConfig(config)
  if err != nil {
    panic(err.Error())
  }

  // templateFile, err := ioutil.ReadFile(*templatePath)
  // if err != nil {
  //   panic(err.Error())
  // }

  // &Service{
  //   ObjectMeta:k8s_io_apimachinery_pkg_apis_meta_v1.ObjectMeta{
  //     Name:sshd-service,
  //     GenerateName:,
  //     Namespace:default,
  //     SelfLink:/api/v1/namespaces/default/services/sshd-service,
  //     UID:6e717261-c1dd-11e7-b876-525400615db6,
  //     ResourceVersion:196871,
  //     Generation:0,
  //     CreationTimestamp:2017-11-04 20:57:24 -0700 PDT,
  //     DeletionTimestamp:<nil>,
  //     DeletionGracePeriodSeconds:nil,
  //     Labels:map[string]string{name: sshd-pod,},
  //     Annotations:map[string]string{},
  //     OwnerReferences:[],
  //     Finalizers:[],
  //     ClusterName:,
  //     Initializers:nil,
  //   },
  //   Spec:ServiceSpec{
  //     Ports:[
  //       { TCP 2222 {0 2222 } 32222}
  //     ],
  //     Selector:map[string]string{k8s-app: sshd,},
  //     ClusterIP:10.3.0.153,
  //     Type:NodePort,
  //     ExternalIPs:[],
  //     SessionAffinity:None,
  //     LoadBalancerIP:,
  //     LoadBalancerSourceRanges:[],
  //     ExternalName:,
  //     ExternalTrafficPolicy:Cluster,
  //     HealthCheckNodePort:0,
  //     PublishNotReadyAddresses:false,
  //     SessionAffinityConfig:nil,
  //   },
  //   Status:ServiceStatus{
  //     LoadBalancer:LoadBalancerStatus{
  //       Ingress:[],
  //     },
  //   },
  // }


  tmpl := template.New("template")
  servicesMap = make(map[Key]PortMap)

  watchlist := cache.NewListWatchFromClient(clientset.Core().RESTClient(), "services", apiv1.NamespaceDefault,
    fields.Everything())
  _, controller := cache.NewInformer(
    watchlist,
    &apiv1.Service{},
    time.Second * 0,
    cache.ResourceEventHandlerFuncs{

      AddFunc: func(obj interface{}) {
        service := obj.(*apiv1.Service)

        for _, value := range service.Spec.Ports {
          servicesMap[Key{service.Name, value.Name}] = PortMap{value.NodePort, value.Port}
        }

        f, _ := os.OpenFile(*outPath, os.O_CREATE|os.O_WRONLY, 0777)
        defer f.Close()

        tmpl, _ = tmpl.ParseFiles(*templatePath)
        tmpl.Execute(f, servicesMap)
      },

      DeleteFunc: func(obj interface{}) {
        service := obj.(*apiv1.Service)

        for _, value := range service.Spec.Ports {
          delete(servicesMap, Key{service.Name, value.Name})
        }

        f, _ := os.OpenFile(*outPath, os.O_CREATE|os.O_WRONLY, 0777)
        defer f.Close()

        tmpl, _ = tmpl.ParseFiles(*templatePath)
        tmpl.Execute(f, servicesMap)
      },

      UpdateFunc:func(oldObj, newObj interface{}) {
        service := newObj.(*apiv1.Service)

        for _, value := range service.Spec.Ports {
          servicesMap[Key{service.Name, value.Name}] = PortMap{value.NodePort, value.Port}
        }

        f, _ := os.OpenFile(*outPath, os.O_CREATE|os.O_WRONLY, 0777)
        defer f.Close()

        tmpl, _ = tmpl.ParseFiles(*templatePath)
        tmpl.Execute(f, servicesMap)
      },
    },
  )
  stop := make(chan struct{})
  go controller.Run(stop)
  for {
    time.Sleep(time.Second)
  }
}
