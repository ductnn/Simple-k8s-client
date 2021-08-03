package main

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func getPodObject() *core.Pod {
	return &core.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "app",
			Namespace: "default",
			Labels: map[string]string{
				"app": "app",
			},
		},
		Spec: core.PodSpec{
			Containers: []core.Container{
				{
					Name:            "green-rain",
					Image:           "ductn4/green-rain",
					ImagePullPolicy: core.PullIfNotPresent,
					Ports: []core.ContainerPort{
						{
							Name:          "container-port",
							ContainerPort: 8000,
							Protocol:      core.ProtocolTCP,
						},
					},
				},
			},
		},
	}
}

func getServicePod() *core.Service {
	return &core.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "service-app",
			Namespace: "default",
		},
		Spec: core.ServiceSpec{
			Selector: map[string]string{
				"app": "app",
			},
			Type: core.ServiceTypeNodePort,
			Ports: []core.ServicePort{
				{
					Name: "service-port",
					Port: 8000,
					TargetPort: intstr.IntOrString{
						Type:   intstr.String,
						StrVal: "container-port",
					},
				},
			},
		},
	}
}

func main() {
	rules := clientcmd.NewDefaultClientConfigLoadingRules()
	kubeconfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		rules,
		&clientcmd.ConfigOverrides{},
	)

	config, err := kubeconfig.ClientConfig()

	if err != nil {
		panic(err)
	}

	clientset := kubernetes.NewForConfigOrDie(config)

	// nodeList, err := clientset.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})

	if err != nil {
		panic(err)
	}

	var numberOfReplica int32
	numberOfReplica = 2

	deploymentsClient := clientset.AppsV1().Deployments(core.NamespaceDefault)

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "app",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &numberOfReplica,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "app",
				},
			},
			Template: core.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "app",
					},
				},
				Spec: core.PodSpec{
					Containers: []core.Container{
						{
							Name:            "green-rain",
							Image:           "ductn4/green-rain",
							ImagePullPolicy: core.PullIfNotPresent,
							Ports: []core.ContainerPort{
								{
									Name:          "container-port",
									ContainerPort: 8000,
									Protocol:      core.ProtocolTCP,
								},
							},
						},
					},
				},
			},
		},
	}

	// Create Deployment
	fmt.Println("Creating deployment...")
	result, err := deploymentsClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())

	// for _, n := range nodeList.Items {
	// 	fmt.Println(n.Name)
	// }

	// serviceClient := clientset.AppsV1().Services(apiv1.NamespaceDefault)

	// newService := &corev1.Service{
	// 	ObjectMeta: metav1.ObjectMeta{
	// 		Name: "serviceapp",
	// 		Namespace: "default",
	// 		Labels: map[string]string{
	// 			"app": "myapp",
	// 		},
	// 	},
	// 	Spec: corev1.ServiceSpec{
	// 		Ports:    nil,
	// 		Selector: nil,
	// 	},
	// }

	// // Create Service
	// fmt.Println("Creating service...")
	// // result, err := servicesClient.Create(service)
	// service, err := clientset.CoreV1().Services("default").Create(context.Background(), newService, metav1.CreateOptions{})
	// if err != nil {
	// 	panic(err)
	// }

	// Deployment
	// newPod := &core.Pod{
	// 	ObjectMeta: metav1.ObjectMeta{
	// 		Name: "test-pod",
	// 	},
	// 	Spec: core.PodSpec{
	// 		Containers: []core.Container{
	// 			{
	// 				Name:  "greenrain",
	// 				Image: "ductn4/green-rain:latest",
	// 				// Command: []string{"sleep", "10000"},
	// 			},
	// 		},
	// 	},
	// }

	svc := getServicePod()
	svc, err = clientset.CoreV1().Services(svc.Namespace).Create(context.Background(), svc, metav1.CreateOptions{})

	if err != nil {
		panic(err)
	}

	// pod := getPodObject()
	// pod, err = clientset.CoreV1().Pods(pod.Namespace).Create(context.Background(), pod, metav1.CreateOptions{})

	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println(pod)
}
