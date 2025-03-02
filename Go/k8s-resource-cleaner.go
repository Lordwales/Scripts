package main

import (
	"context"
	"fmt"
	"log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main() {
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatalf("Error getting cluster config: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Error creating Kubernetes client: %v", err)
	}

	cleanupPods(clientset)
	cleanupServices(clientset)
}

func cleanupPods(clientset *kubernetes.Clientset) {
	pods, err := clientset.CoreV1().Pods("default").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Error listing pods: %v", err)
	}

	for _, pod := range pods.Items {
		if pod.Status.Phase == "Succeeded" || pod.Status.Phase == "Failed" {
			fmt.Printf("Deleting pod: %s\n", pod.Name)
			err := clientset.CoreV1().Pods("default").Delete(context.TODO(), pod.Name, metav1.DeleteOptions{})
			if err != nil {
				log.Printf("Error deleting pod %s: %v", pod.Name, err)
			}
		}
	}
}

func cleanupServices(clientset *kubernetes.Clientset) {
	services, err := clientset.CoreV1().Services("default").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Error listing services: %v", err)
	}

	for _, svc := range services.Items {
		if len(svc.Spec.Selector) == 0 {
			fmt.Printf("Deleting unused service: %s\n", svc.Name)
			err := clientset.CoreV1().Services("default").Delete(context.TODO(), svc.Name, metav1.DeleteOptions{})
			if err != nil {
				log.Printf("Error deleting service %s: %v", svc.Name, err)
			}
		}
	}
}
