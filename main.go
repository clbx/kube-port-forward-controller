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
			AddFunc: func(obj interface{}) {

				service := obj.(*v1.Service)
				for key, value := range service.Annotations {
					if key == "kube-router-port-forward/ports" {
						fmt.Printf("Found new service with port annotation: %s value: %s\n", service.Name, value)
					}
				}

			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				// Check if the old service had an annotation
				oldService := oldObj.(*v1.Service)
				for key, value := range oldService.Annotations {
					if key == "kube-router-port-forward/ports" {
						// Remove old port
						fmt.Printf("MOCK: remove old ports: %s\n", value)

						// Add new port
						newService := newObj.(*v1.Service)
						for key, value := range newService.Annotations {
							if key == "kube-router-port-forward/ports" {
								//Add new port
								fmt.Printf("MOCK: add new ports: %s\n", value)
							}
						}
					}
				}

			},
			DeleteFunc: func(obj interface{}) {
				service := obj.(*v1.Service)
				for key, value := range service.Annotations {
					if key == "kube-router-port-forward/ports" {
						fmt.Printf("MOCK: remove port: %s\n", value)
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
