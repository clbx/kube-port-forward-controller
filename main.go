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
			// I don't think this will ever realisctially be hit since when a LoadBalancer is made, it doesn't have an IP
			// and then the update function will catch it when it gets assinged it.
			AddFunc: func(obj interface{}) {
				service := obj.(*v1.Service)
				//fmt.Printf("New Service Added: %s\n", service.Name)
				// Check if the Service is a Load Balancer
				if service.Spec.Type != "LoadBalancer" {
					//fmt.Printf("Service %s is not a Load Balancer.. Skipping\n", service.Name)
					return
				}
				fmt.Printf("Load Balancer %s found\n", service.Name)
				for key, value := range service.Annotations {
					if key == "kube-router-port-forward/ports" {
						// Shitty, but it should assign an IP by this long, I'll figure it out later
						if len(service.Status.LoadBalancer.Ingress) > 0 && service.Status.LoadBalancer.Ingress[0].IP != "" {
							ip := service.Status.LoadBalancer.Ingress[0].IP
							fmt.Printf("MOCK: update add new ip: %s ports: %s\n", ip, value)
						} else {
							fmt.Printf("New Service %s found, but LoadBalancer does not have an IP yet. Ignoring\n", service.Name)
						}
						return
					}
				}
			},

			// When updating a service
			UpdateFunc: func(oldObj, newObj interface{}) {

				oldService := oldObj.(*v1.Service)
				newService := newObj.(*v1.Service)

				// Check if the service is a LoadBalancer
				if newService.Spec.Type != "LoadBalancer" {
					//fmt.Printf("Service %s is not a Load Balancer.. Skipping\n", service.Name)
					return
				}

				fmt.Printf("Update Detected %s\n", oldService.Name)

				oldPort, oldExists := oldService.Annotations["kube-router-port-forward/ports"]
				newPort, newExists := newService.Annotations["kube-router-port-forward/ports"]

				// If the port was removed, remove it from the router
				if oldExists && !newExists {
					fmt.Printf("MOCK: Remove port %s from router", oldPort)
					return
				}

				//If there's just a new port
				if !oldExists && newExists {
					fmt.Printf("MOCK: Add port %s to router", newPort)
					return
				}

				// If the old service and new service have a port
				if oldExists && newExists {
					if oldPort != newPort {
						fmt.Printf("MOCK: Remove ip: %s port: %s from router", getLBIP(oldService), oldPort)
						fmt.Printf("MOCK: Add ip: %s port: %s to router", getLBIP(newService), newPort)
						// If they're the same, nothing needs to be done UNLESS its new, so we add!  If its already there, the adding function will take care of that.
					} else {
						fmt.Printf("MOCK: Add ip: %s port: %s to router", getLBIP(newService), newPort)
					}
					return
				}

			},
			// When deleting a service
			DeleteFunc: func(obj interface{}) {
				service := obj.(*v1.Service)

				//First check if its a LoadBalancer
				if service.Spec.Type != "LoadBalancer" {
					//fmt.Printf("Service %s is not a Load Balancer.. Skipping\n", service.Name)
					return
				}

				if port, exists := service.Annotations["kube-router-port-forward/ports"]; exists {
					fmt.Printf("MOCK: Remove ip: %s port: %s from router", getLBIP(service), port)
				}

				// for key, value := range service.Annotations {
				// 	if key == "kube-router-port-forward/ports" {
				// 		if len(service.Status.LoadBalancer.Ingress) > 0 {
				// 			ip := service.Status.LoadBalancer.Ingress[0].IP
				// 			fmt.Printf("MOCK: remove ip: %s ports: %s\n", ip, value)
				// 		}
				// 	}
				// }
			},
		},
	)

	stop := make(chan struct{})
	defer close(stop)
	go controller.Run(stop)

	// Wait forever
	select {}
}

func getLBIP(service *v1.Service) string {
	if len(service.Status.LoadBalancer.Ingress) > 0 {
		ip := service.Status.LoadBalancer.Ingress[0].IP
		if len(service.Status.LoadBalancer.Ingress) > 1 {
			fmt.Printf("Service %s has more than one ip. Proceeding with %s", service.Name, ip)
		}
		return ip
	} else {
		fmt.Printf("Service %s has no IP", service.Name)
		return ""
	}

}
