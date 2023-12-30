package main

import (
	"fmt"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

func main() {

	// load in cluster config from service account
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	watchlist := cache.NewListWatchFromClient(
		clientset.CoreV1().RESTClient(),
		"services",
		metav1.NamespaceAll,
		fields.Everything(),
	)

	_, controller := cache.NewInformer(
		watchlist,
		&v1.Service{},
		time.Second*0,
		cache.ResourceEventHandlerFuncs{

			// When adding a service
			AddFunc: func(obj interface{}) {
				service := obj.(*v1.Service)
				fmt.Printf("New Service Added: %s\n", service.Name)
				for key, value := range service.Annotations {
					if key == "kube-router-port-forward/ports" {
						if len(service.Status.LoadBalancer.Ingress) > 0 {
							ip := service.Status.LoadBalancer.Ingress[0].IP
							fmt.Printf("MOCK: add new ip: %s, ports: %s\n", ip, value)
							return
						} else {
							fmt.Printf("Add Service %s found, but either isn't a load balancer or doesn't have an IP. Ignoring\n", service.Name)
							return
						}

					}
				}

			},

			// When updating a service
			UpdateFunc: func(oldObj, newObj interface{}) {
				// Check if the old service had an annotation
				oldService := oldObj.(*v1.Service)
				for key, value := range oldService.Annotations {
					if key == "kube-router-port-forward/ports" {
						// Remove old port
						if len(oldService.Status.LoadBalancer.Ingress) > 0 {
							ip := oldService.Status.LoadBalancer.Ingress[0].IP
							fmt.Printf("MOCK: remove ip: %s ports: %s\n", ip, value)
						}
						// Add new port
						newService := newObj.(*v1.Service)
						for key, value := range newService.Annotations {
							if key == "kube-router-port-forward/ports" {
								//Add new port
								if len(newService.Status.LoadBalancer.Ingress) > 0 {
									ip := newService.Status.LoadBalancer.Ingress[0].IP
									fmt.Printf("MOCK: update add new ip: %s ports: %s\n", ip, value)
									return
								} else {
									fmt.Printf("Update Service %s found, but either isn't a load balancer or doesn't have an IP. Ignoring\n", newService.Name)
									return
								}

							}
						}
					}
				}

			},

			// When deleting a service
			DeleteFunc: func(obj interface{}) {
				service := obj.(*v1.Service)
				for key, value := range service.Annotations {
					if key == "kube-router-port-forward/ports" {
						if len(service.Status.LoadBalancer.Ingress) > 0 {
							ip := service.Status.LoadBalancer.Ingress[0].IP
							fmt.Printf("MOCK: remove ip: %s ports: %s\n", ip, value)
						}
					}
				}
			},
		},
	)

	stop := make(chan struct{})
	defer close(stop)
	go controller.Run(stop)

	// Wait forever
	select {}
}
