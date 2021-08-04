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

func getDeploymentObject() *appsv1.Deployment {
	var numberOfReplica int32
	numberOfReplica = 1
	return &appsv1.Deployment{
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
							Name:            "url-shorter",
							Image:           "ductn4/urlshorter:1.0.0",
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
}

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
						IntVal: 8000,
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

	if err != nil {
		panic(err)
	}

	deploymentsClient := clientset.AppsV1().Deployments(core.NamespaceDefault)

	if err != nil {
		panic(err)
	}

	// Create Deployment
	deployment := getDeploymentObject()
	fmt.Println("Creating deployment...")
	result, err := deploymentsClient.Create(context.TODO(), deployment, metav1.CreateOptions{})

	if err != nil {
		panic(err)
	}
	fmt.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())

	// Create Service
	svc := getServicePod()
	fmt.Println("Creating service...")
	svc, err = clientset.CoreV1().Services(svc.Namespace).Create(context.Background(), svc, metav1.CreateOptions{})

	if err != nil {
		panic(err)
	}
	fmt.Printf("Created service %q.\n", result.GetObjectMeta().GetName())
}
